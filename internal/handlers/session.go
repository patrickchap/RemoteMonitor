package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) NewSession(co CookieOpts, authFailed echo.HandlerFunc) {

	var store *sessions.CookieStore
	var mgr *manager

	store = sessions.NewCookieStore([]byte(co.Secret))
	sess := sessions.NewSession(store, co.Name)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   co.MaxAge,
		HttpOnly: true,
		Secure:   true,
	}
	mgr = &manager{
		session:    sess,
		cookie:     co,
		authFailed: authFailed,
	}
	h.Manager = mgr
}

func (h *Handler) GetSession(c echo.Context) *sessions.Session {

	var sess = h.Manager.session
	if s, err := session.Get(h.Manager.cookie.Name, c); err == nil {
		sess = s
	}
	return sess
}

func (h *Handler) Set(c echo.Context, key string, value interface{}) error {
	sess := h.GetSession(c)
	sess.Values[key] = value

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Get(c echo.Context, key string) interface{} {
	sess := h.GetSession(c)
	return sess.Values[key]
}
