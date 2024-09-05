package handlers

import (
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	"database/sql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(c echo.Context) error {
	var isLoggedIn bool
	isLoggedIn = false
	if isLoggedIn {
		return c.Redirect(302, "/dashboard")
	}

	return helpers.RenderTemplate(c, views.Login())
}

type PostLoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (h *Handler) PostLogin(c echo.Context) error {

	//validate form data
	req := new(PostLoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	user, err := h.Store.GetUserForAuth(ctx, sql.NullString{String: req.Email, Valid: true})
	if err != nil {
		return helpers.RenderTemplate(c, views.LoginForm([]string{"User not Found"}))
	}

	checkPassord := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(req.Password))
	if checkPassord != nil {
		return helpers.RenderTemplate(c, views.LoginForm([]string{"Invalid Password"}))
	}

	h.Set(c, "userId", user.ID)
	h.Set(c, "userEmail", user.Email)
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return nil
}
