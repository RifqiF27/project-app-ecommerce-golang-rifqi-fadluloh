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
func (h *AuthHandler) GetDetailUserHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	users, err := h.authService.GetDetailUserService(userID)
	if err != nil {
		h.Log.Error("Handler: Error getting users", zap.Error(err))
		h.Log.Debug("Handler: Error getting users", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", users)
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

func (h *AuthHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request to update user", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	var input model.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Log.Error("Handler: Failed to decode request body", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if len(input.Address) == 0 || (len(input.Address) == 1 && input.Address[0] == "") {
		h.Log.Warn("Handler: Validation failed for address field", zap.Any("address", input.Address))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Address cannot be empty", nil)
		return
	}
	if err := h.validator.ValidateStruct(input); err != nil {
		formattedError := helper.FormatValidationError(err)
		h.Log.Error("Handler: validation failed", zap.String("error", formattedError))
		h.Log.Debug("Handler: validation failed", zap.String("error", formattedError))
		helper.SendJSONResponse(w, http.StatusBadRequest, formattedError, nil)
		return
	}
	updatedUser, err := h.authService.UpdateUserService(userID, input.Name, input.Email, input.Phone, input.Password, input.Address)
	if err != nil {
		h.Log.Error("Handler: Failed to update user", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, "Failed to update user", nil)
		return
	}

	h.Log.Info("Handler: User updated successfully", zap.Int("userID", updatedUser.ID))
	helper.SendJSONResponse(w, http.StatusOK, "User updated successfully", updatedUser)
}

func (h *AuthHandler) CreateAddressHandler(w http.ResponseWriter, r *http.Request) {
    h.Log.Info("Handler: Received request to create address", zap.String("method", r.Method), zap.String("path", r.URL.Path))

    // Ambil userID dari context
    userID, ok := r.Context().Value("userID").(int)
    if !ok {
        h.Log.Error("Handler: userID not found in context")
        helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
        return
    }

    // Parse body JSON
    var input struct {
        NewAddress string `json:"new_address"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        h.Log.Error("Handler: Failed to decode request body", zap.Error(err))
        helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    // Validasi input
    if input.NewAddress == "" {
        h.Log.Warn("Handler: Address cannot be empty", zap.Any("input", input))
        helper.SendJSONResponse(w, http.StatusBadRequest, "Address cannot be empty", nil)
        return
    }

    // Panggil service untuk menambahkan alamat baru
    updatedUser, err := h.authService.CreateAddressService(userID, input.NewAddress)
    if err != nil {
        h.Log.Error("Handler: Failed to create address", zap.Error(err))
        helper.SendJSONResponse(w, http.StatusInternalServerError, "Failed to create address", nil)
        return
    }

    h.Log.Info("Handler: Address created successfully", zap.Int("userID", updatedUser.ID))
    helper.SendJSONResponse(w, http.StatusOK, "Address created successfully", updatedUser)
}
