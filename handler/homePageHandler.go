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

type HomePageHandler struct {
	service   service.HomePageService
	Log       *zap.Logger
	validator *helper.Validator
	config    util.Configuration
}

func NewHomePageHandler(service service.HomePageService, logger *zap.Logger, config util.Configuration) *HomePageHandler {
	return &HomePageHandler{service: service, Log: logger, validator: helper.NewValidator(), config: config}
}

func (h *HomePageHandler) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	name := r.URL.Query().Get("name")
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit == 0 {
		limit = 5
	}
	if page == 0 {
		page = 1
	}

	products, totalItems, totalPages, err := h.service.GetAllProductsService(name, categoryID, limit, page)
	if err != nil {
		h.Log.Error("Handler: Error getting products", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(products) == 0 {
		h.Log.Warn("Handler: No products found", zap.String("name", name), zap.Int("category_id", categoryID))
		helper.SendJSONResponse(w, http.StatusNotFound, "No products found", nil)
		return
	}

	helper.SendJSONResponsePagination(w, page, limit, totalItems, totalPages, http.StatusOK, "", products)
}

func (h *HomePageHandler) GetByIdProductHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	item, err := h.service.GetByIdProductService(id)
	if err != nil {
		h.Log.Error("Handler: Product not found", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusNotFound, "Product not found", nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", item)
}
func (h *HomePageHandler) GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	item, err := h.service.GetAllCategoriesService()
	if err != nil {
		h.Log.Error("Handler: Error getting categories", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", item)
}

func (h *HomePageHandler) GetAllBannersHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	banners, err := h.service.GetAllBannersService()
	if err != nil {
		h.Log.Error("Handler: Error getting banners", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", banners)
}
func (h *HomePageHandler) GetAllBestSellingProductsHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit == 0 {
		limit = 5
	}
	if page == 0 {
		page = 1
	}

	products, totalItems, totalPages, err := h.service.GetAllBestSellingProductsService(limit, page)
	if err != nil {
		h.Log.Error("Handler: Error getting products best selling", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if len(products) == 0 {
		h.Log.Warn("Handler: No products found", zap.Int("page", page), zap.Int("limit", limit))
		helper.SendJSONResponse(w, http.StatusNotFound, "No products found", nil)
		return
	}

	helper.SendJSONResponsePagination(w, page, limit, totalItems, totalPages, http.StatusOK, "", products)
}

func (h *HomePageHandler) GetAllWeeklyPromotionProductsHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	products, err := h.service.GetAllWeeklyPromotionProductsService()
	if err != nil {
		h.Log.Error("Handler: Error getting products weekly promotion", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", products)
}

func (h *HomePageHandler) GetAllRecommentsProductsHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Handler: Received request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

	products, err := h.service.GetAllRecommentsProductsService()
	if err != nil {
		h.Log.Error("Handler: Error getting products recomments", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "", products)
}

func (h *HomePageHandler) AddWishlistHandler(w http.ResponseWriter, r *http.Request) {
	var wishlist model.Wishlist

	if err := json.NewDecoder(r.Body).Decode(&wishlist); err != nil {
		h.Log.Error("Handler: invalid request payload", zap.Error(err))
		h.Log.Debug("Handler: invalid request payload", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	wishlist.UserID = userID

	err := h.service.AddWishlistService(wishlist)
	if err != nil {
		h.Log.Error("Handler: add wishlist failed", zap.Error(err))
		h.Log.Debug("Handler: add wishlist failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusCreated, "Wishlist successfully added", nil)

}

func (h *HomePageHandler) DeleteWishlistHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Log.Error("Handler: Invalid wishlist ID", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.Log.Error("Handler: userID not found in context")
		helper.SendJSONResponse(w, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	err = h.service.DeleteWishlistService(id, userID)
	if err != nil {
		h.Log.Error("Handler: delete wishlist failed", zap.Error(err))
		helper.SendJSONResponse(w, http.StatusConflict, err.Error(), nil)
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "Wishlist successfully deleted", nil)
}
