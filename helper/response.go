package helper

import (
	"ecommerce/model"
	"encoding/json"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := model.Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func SendJSONResponsePagination(w http.ResponseWriter, page, limit, totalItems, totalPages, status int, message string, data interface{}) {
	response := model.Response{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Status:     status,
		Message:    message,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}