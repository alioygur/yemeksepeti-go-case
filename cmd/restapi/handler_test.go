package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alioygur/yemekepeti-go-case/herodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_get(t *testing.T) {
	mockStore := new(herodb.MockStore)
	handler := &Handler{Storage: mockStore}

	t.Run("normal", func(t *testing.T) {
		mockStore.On("Get", context.Background(), "name").Return("ali", nil).Once()

		w := httptest.NewRecorder()

		handler.get(w, httptest.NewRequest(http.MethodGet, "/?key=name", nil))

		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
		assert.Equal(t, `{"value":"ali"}`, strings.TrimSuffix(string(body), "\n"))

		mockStore.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockStore.On("Get", context.Background(), "name").Return("", herodb.ErrNotfound).Once()

		w := httptest.NewRecorder()

		handler.get(w, httptest.NewRequest(http.MethodGet, "/?key=name", nil))

		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
		assert.Equal(t, `{"error":"the key is not found"}`, strings.TrimSuffix(string(body), "\n"))

		mockStore.AssertExpectations(t)
	})

	t.Run("server error", func(t *testing.T) {
		mockStore.On("Get", context.Background(), "name").Return("", errors.New("server error")).Once()

		w := httptest.NewRecorder()

		handler.get(w, httptest.NewRequest(http.MethodGet, "/get?key=name", nil))

		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
		assert.Equal(t, `{"error":"server error"}`, strings.TrimSuffix(string(body), "\n"))

		mockStore.AssertExpectations(t)
	})
}

func TestHandler_set(t *testing.T) {
	mockStore := new(herodb.MockStore)
	handler := &Handler{Storage: mockStore}

	t.Run("normal", func(t *testing.T) {
		mockStore.On("Set", context.Background(), "name", "ali").Return(nil).Once()

		w := httptest.NewRecorder()

		handler.set(w, httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(`{"key": "name", "value": "ali"}`)))

		res := w.Result()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		mockStore.AssertExpectations(t)
	})

	t.Run("missing values", func(t *testing.T) {
		w := httptest.NewRecorder()

		handler.set(w, httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(`{"value": "ali"}`)))

		res := w.Result()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		mockStore.AssertExpectations(t)
	})
}
