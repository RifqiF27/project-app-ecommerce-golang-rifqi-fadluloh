package handler

import (
	"ecommerce/helper"
	"ecommerce/model"
	"ecommerce/service"
	"ecommerce/util"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	Log         *zap.Logger
	validator   *helper.Validator
	config      util.Configuration
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger, config util.Configuration) *AuthHandler {
	return &AuthHandler{authService: authService, Log: logger, validator: helper.NewValidator(), config: config}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		h.Log.Debug("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.validator.ValidateStruct(req); err != nil {
		formattedError := helper.FormatValidationError(err)
		h.Log.Error("Handler: validation failed", zap.String("error", formattedError))
		h.Log.Debug("Handler: validation failed", zap.String("error", formattedError))
		helper.SendJSONResponse(w, http.StatusBadRequest, formattedError, nil)
		return
	}

	err := h.authService.RegisterService(req)
	if err != nil {
		h.Log.Error("Handler: register failed", zap.Error(err))
		h.Log.Debug("Handler: register failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "Registration successful", nil)

}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		h.Log.Debug("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.validator.ValidateLoginStruct(req); err != nil {
		formattedError := helper.FormatValidationError(err)
		h.Log.Error("Handler: validation failed", zap.String("error", formattedError))
		h.Log.Debug("Handler: validation failed", zap.String("error", formattedError))
		helper.SendJSONResponse(w, http.StatusBadRequest, formattedError, nil)
		return
	}

	users, err := h.authService.LoginService(req)
	if err != nil {
		h.Log.Error("Handler: login failed", zap.Error(err))
		h.Log.Debug("Handler: login failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	token, err := helper.GenerateToken()
	if err != nil {
		return
	}
	session := &model.Session{
		UserID:    users.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.authService.CreateSessionService(*session); err != nil {
		return
	}

	users.Token = token

	helper.SendJSONResponse(w, http.StatusOK, "login success", users)
}

func (h *AuthHandler) GetAllAddressHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}
	
	address, err := h.authService.GetAllAddressService(userID)
	if err != nil {
		h.Log.Error("Handler: Error getting address", zap.Error(err))
		h.Log.Debug("Handler: Error getting address", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", address)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if token == "" {
		h.Log.Error("Handler: token required", zap.String("error: ", token))
		h.Log.Debug("Handler: token required", zap.String("error: ", token))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "Token required", nil)
		return
	}

	_, err := h.authService.VerifyToken(token)
	if err != nil {
		h.Log.Error("Handler: token invalid or expired", zap.Error(err))
		h.Log.Debug("Handler: token invalid or expired", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
		return
	}

	if err := h.authService.Logout(token); err != nil {
		h.Log.Error("Handler: logout failed", zap.Error(err))
		h.Log.Debug("Handler: logout failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, "Logout failed", nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "Logout success", nil)
}
