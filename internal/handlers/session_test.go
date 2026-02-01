package handlers_test

import (
	"RemoteMonitor/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func setupTestHandler() *handlers.Handler {
	h := &handlers.Handler{}
	co := handlers.CookieOpts{
		Name:   "test_session",
		Secret: "test-secret-key-32-bytes-long!!",
		MaxAge: 3600,
	}
	h.NewSession(co, func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "unauthorized")
	})
	return h
}

func setupTestContext(h *handlers.Handler) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	store := sessions.NewCookieStore([]byte("test-secret-key-32-bytes-long!!"))
	e.Use(session.Middleware(store))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestHandler_Get(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue any
		want     any
	}{
		{
			name:     "get existing string value",
			key:      "username",
			setValue: "testuser",
			want:     "testuser",
		},
		{
			name:     "get existing int value",
			key:      "user_id",
			setValue: 123,
			want:     123,
		},
		{
			name:     "get non-existent key returns nil",
			key:      "nonexistent",
			setValue: nil,
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupTestHandler()
			c, _ := setupTestContext(h)

			if tt.setValue != nil {
				err := h.Set(c, tt.key, tt.setValue)
				if err != nil {
					t.Fatalf("failed to set value: %v", err)
				}
			}

			got := h.Get(c, tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   any
		wantErr bool
	}{
		{
			name:    "set string value",
			key:     "username",
			value:   "testuser",
			wantErr: false,
		},
		{
			name:    "set int value",
			key:     "user_id",
			value:   42,
			wantErr: false,
		},
		{
			name:    "set bool value",
			key:     "is_admin",
			value:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupTestHandler()
			c, _ := setupTestContext(h)

			gotErr := h.Set(c, tt.key, tt.value)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Set() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Set() succeeded unexpectedly")
			}

			got := h.Get(c, tt.key)
			if got != tt.value {
				t.Errorf("Set() value not persisted, got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestHandler_GetSession(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "get session returns non-nil session",
		},
		{
			name: "get session returns session with values map",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupTestHandler()
			c, _ := setupTestContext(h)

			got := h.GetSession(c)
			if got == nil {
				t.Error("GetSession() = nil, want non-nil session")
			}
			if got.Values == nil {
				t.Error("GetSession().Values = nil, want non-nil map")
			}
		})
	}
}

func TestHandler_NewSession(t *testing.T) {
	tests := []struct {
		name   string
		co     handlers.CookieOpts
		wantOK bool
	}{
		{
			name: "create session with valid options",
			co: handlers.CookieOpts{
				Name:   "my_session",
				Secret: "secret-key-that-is-32-bytes-lng",
				MaxAge: 7200,
			},
			wantOK: true,
		},
		{
			name: "create session with empty name",
			co: handlers.CookieOpts{
				Name:   "",
				Secret: "secret-key",
				MaxAge: 3600,
			},
			wantOK: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handlers.Handler{}
			authFailed := func(c echo.Context) error {
				return c.String(http.StatusUnauthorized, "unauthorized")
			}

			h.NewSession(tt.co, authFailed)

			if tt.wantOK && h.Manager == nil {
				t.Error("NewSession() did not set Manager")
			}
		})
	}
}
