package tg

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_SendMessage(t *testing.T) {
	const token = "tok"
	t.Setenv("TELEGRAM_TOKEN", token)

	var called bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		require.Equal(t, "/bot"+token+"/sendMessage", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(srv.Client())
	c.apiURL = srv.URL

	err := c.SendMessage(t.Context(), "1", "hello")
	require.NoError(t, err)
	require.True(t, called)
}

func TestClient_SendMessage_Error(t *testing.T) {
	const token = "tok"
	t.Setenv("TELEGRAM_TOKEN", token)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad"))
	}))
	defer srv.Close()

	c := New(srv.Client())
	c.apiURL = srv.URL

	err := c.SendMessage(t.Context(), "1", "hello")
	require.Error(t, err)
}
