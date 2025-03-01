package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/models"
	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/scripts"
	"github.com/Dyrakavalyanie/Clothes_shop/services/catalog/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct {
	ConnPool *pgxpool.Pool
	Cfg config.Config
}

func (handler *Handler)AddItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var item models.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var itemID *models.ItemID
	itemID, err = scripts.AddItem(ctx, handler.ConnPool, &item)

	if err != nil {
		http.Error(w, "Adding error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(itemID)
	if err != nil {
		http.Error(w, "Error in encoding response", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler)UpdateItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var item models.ItemUpdate
	err = json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = scripts.UpdateItem(ctx, handler.ConnPool, &item)

	if err != nil {
		http.Error(w, "Updating error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler)DeleteItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var item models.ItemID
	err = json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = scripts.DeleteItem(ctx, handler.ConnPool, &item)

	if err != nil {
		http.Error(w, "Deleting error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler)GetItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	var item models.ItemID
	err = json.Unmarshal(body, &item)
	if err != nil {
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var foundItem *models.Item
	foundItem, err = scripts.GetItem(ctx, handler.ConnPool, &item)
	if err != nil {
		http.Error(w, "Getting error", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(foundItem)
	if err != nil {
		http.Error(w, "Error in encoding response", http.StatusInternalServerError)
		return
	}
}