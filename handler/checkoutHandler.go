package handler

import (
	"ecommerce/helper"
	"ecommerce/model"
	"ecommerce/service"
	"ecommerce/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type CheckoutHandler struct {
	service   service.CheckoutService
	Log       *zap.Logger
	validator *helper.Validator
	config    util.Configuration
}

func NewCheckoutHandler(service service.CheckoutService, logger *zap.Logger, config util.Configuration) *CheckoutHandler {
	return &CheckoutHandler{service: service, Log: logger, validator: helper.NewValidator(), config: config}
}

func (h *CheckoutHandler) GetAllCartHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	carts, err := h.service.GetAllCartService(userID)
	if err != nil {
		h.Log.Error("Handler: Error getting carts", zap.Error(err))
		h.Log.Debug("Handler: Error getting carts", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", carts)
}
func (h *CheckoutHandler) GetTotalCartHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var cart struct {
		Total_carts int `json:"total_carts"`
	}
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	carts, err := h.service.GetTotalCartService(userID)
	if err != nil {
		h.Log.Error("Handler: Error getting total carts", zap.Error(err))
		h.Log.Debug("Handler: Error getting total carts", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	
	cart.Total_carts = carts.TotalCarts

	helper.SendJSONResponse(w, http.StatusOK, "", cart)
}
func (h *CheckoutHandler) AddCartHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var cart model.Checkout

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		h.Log.Debug("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	cart.UserID = userID

	err := h.service.AddCartService(cart)
	if err != nil {
		h.Log.Error("Handler: add cart failed", zap.Error(err))
		h.Log.Debug("Handler: add cart failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusCreated, "Cart successfully added", nil)
}

func (h *CheckoutHandler) DeleteCartHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Log.Error("Handler: Invalid cart ID", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid cart ID", nil)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	err = h.service.DeleteCartService(id, userID)
	if err != nil {
		h.Log.Error("Handler: delete cart failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "Cart successfully deleted", nil)
}

func (h *CheckoutHandler) UpdateCartHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Log.Error("Handler: Invalid cart ID", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid cart ID", nil)
		return
	}

	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	var cart *model.Checkout

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		h.Log.Debug("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		h.Log.Debug("Handler: userID not found in context", zap.Int("userID", userID))
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	cart.ID = id
	cart.UserID = userID

	cart, err = h.service.UpdateCartService(cart.UserID, cart.ProductID, cart.Quantity)
	if err != nil {
		if err.Error() == "cart not found" {
			h.Log.Warn("Handler: No cart found for user", zap.Int("userID", cart.UserID), zap.Int("productID", cart.ProductID))
			helper.SendJSONResponse(w, http.StatusNotFound, "Cart item not found", nil)
			return
		}
		h.Log.Error("Handler: updated cart failed", zap.Error(err))
		h.Log.Debug("Handler: updated cart failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if cart == nil {

		helper.SendJSONResponse(w, http.StatusOK, "Cart item deleted", nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "Cart successfully updated", cart)
}

func (h *CheckoutHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	var requestData struct {
		ProductID    []int `json:"product_id"`
		AddressIndex int   `json:"address_index"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	h.Log.Info("Handler: Decoded request data", zap.Int("address_index", requestData.AddressIndex))

	if len(requestData.ProductID) == 0 {
		h.Log.Warn("Handler: No items in cart")
		helper.SendJSONResponse(w, http.StatusBadRequest, "Cart is empty", nil)
		return
	}

	orderResponse, err := h.service.CreateOrderService(userID, requestData.ProductID, requestData.AddressIndex)

	if err != nil {
		h.Log.Error("Handler: failed to create order", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, "Failed to create order", nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusCreated, "Order successfully created", orderResponse)
}
