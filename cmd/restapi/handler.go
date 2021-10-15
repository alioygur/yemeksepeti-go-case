package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alioygur/yemekepeti-go-case/herodb"
)

type Storer interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Flush(ctx context.Context) error
}

type Handler struct {
	Storage Storer
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	// get and validate the key
	key := r.URL.Query().Get("key")
	if key == "" {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("the key is required"))
		return
	}

	// try to get from storage
	val, err := h.Storage.Get(r.Context(), key)
	if err != nil {
		// is it client error?
		if errors.Is(err, herodb.ErrNotfound) {
			JSONError(w, http.StatusBadRequest, fmt.Errorf("the key is not found"))
			return
		}
		// server error
		JSONError(w, http.StatusInternalServerError, err)
		// ReportError(r, err)
		return
	}

	// everything is ok
	JSON(w, http.StatusOK, map[string]string{"value": val})
}

type setRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handler) set(w http.ResponseWriter, r *http.Request) {
	// bind request payload to object
	var req setRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}

	// validate the request
	if req.Key == "" || req.Value == "" {
		JSONError(w, http.StatusBadRequest, errors.New("key and value are required"))
		return
	}

	// try to store
	if err := h.Storage.Set(r.Context(), req.Key, req.Value); err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		// ReportError(r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) flush(w http.ResponseWriter, r *http.Request) {
	if err := h.Storage.Flush(r.Context()); err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		// ReportError(r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
