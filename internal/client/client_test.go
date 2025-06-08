package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type resp struct {
	Value string `json:"value"`
}

func TestClient_DoJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{Value: "ok"})
	}))
	defer srv.Close()

	c := New(srv.Client().Transport)
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	var out resp
	err = c.DoJSON(t.Context(), req, &out)
	require.NoError(t, err)
	require.Equal(t, "ok", out.Value)
}

func TestClient_DoJSON_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusBadRequest)
	}))
	defer srv.Close()

	c := New(srv.Client().Transport)
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	err = c.DoJSON(t.Context(), req, nil)
	var ue *ErrUnexpectedStatus
	require.ErrorAs(t, err, &ue)
	require.Equal(t, http.StatusBadRequest, ue.Code)
}

func TestClient_DoJSON_BadScheme(t *testing.T) {
	c := New(nil)
	req, err := http.NewRequest(http.MethodGet, "ftp://example.com", nil)
	require.NoError(t, err)
	err = c.DoJSON(t.Context(), req, nil)
	require.ErrorIs(t, err, ErrBadResponseScheme)
}
