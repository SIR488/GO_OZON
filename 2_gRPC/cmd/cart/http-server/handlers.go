package httpserver

import (
	"MyMod/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Service interface {
	InsertSKU(userId, skuId, count uint32) error
	GetCart(userId uint32) (domain.CartResponse, error)
}
type UI struct {
	service Service
}

func NewUI(service Service) *UI {
	ui := &UI{service: service}
	return ui
}

func (ui *UI) PostSKU(w http.ResponseWriter, r *http.Request) {
	strUserId := r.PathValue("user_id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if userId < 0 {
		http.Error(w, "id must be not negative", http.StatusBadRequest)
		return
	}

	strSkuId := r.PathValue("sku_id")
	skuId, err := strconv.Atoi(strSkuId)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if skuId < 0 {
		http.Error(w, "sku_id must be not negative", http.StatusBadRequest)
		return
	}
	var body struct {
		Count uint32 `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	err = ui.service.InsertSKU(uint32(userId), uint32(skuId), body.Count)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, domain.ErrPreconditionFailed):
		http.Error(w, "sku not found", http.StatusPreconditionFailed)
	case errors.Is(err, domain.ErrInvalidArgument): //кейс не сработает, но добавил
		http.Error(w, "invalid argument", http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (ui *UI) GetCart(w http.ResponseWriter, r *http.Request) {
	strUserId := r.PathValue("user_id")
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		http.Error(w, "ошибка парсинга", http.StatusInternalServerError)
		return
	}
	if userId < 0 {
		http.Error(w, "user_id must be not negative", http.StatusBadRequest)
		return
	}
	cart, err := ui.service.GetCart(uint32(userId))
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)

}
