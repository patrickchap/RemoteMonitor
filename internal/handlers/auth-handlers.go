package handlers

import (
	database "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// move to env.......
const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

type Claims struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(c echo.Context) error {

	if c.Get("user") == nil {
		return helpers.RenderTemplate(c, views.Login())
	}
	return c.Redirect(302, "/admin/dashboard")

}

func (h *Handler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "user",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	})

	c.SetCookie(&http.Cookie{
		Name:    accessTokenCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	})

	c.SetCookie(&http.Cookie{
		Name:    refreshTokenCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	})

	c.SetCookie(&http.Cookie{
		Name:    "auth-session",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
		MaxAge:  -1,
	})

	c.Response().Header().Set("HX-Redirect", "/admin/dashboard")
	return nil
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
	user, err := h.Store.GetUserByEmail(ctx, sql.NullString{String: req.Email, Valid: true})
	if err != nil {
		return helpers.RenderTemplate(c, views.LoginForm([]string{"User not Found"}))
	}

	/* checkPassord := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(req.Password))
	if checkPassord != nil {
		return helpers.RenderTemplate(c, views.LoginForm([]string{"Invalid Password"}))
	} */

	err = GenerateTokensAndSetCookies(&user, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	h.Set(c, "userId", user.ID)
	h.Set(c, "userEmail", user.Email)
	if target := c.QueryParam("target"); target != "" {
		c.Response().Header().Set("HX-Redirect", target)
		return nil
	}
	c.Response().Header().Set("HX-Redirect", "/admin/dashboard")
	return nil
}

func GenerateTokensAndSetCookies(user *database.User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)

	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func generateRefreshToken(user *database.User) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *database.User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateToken(user *database.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		Name: user.FirstName.String + " " + user.LastName.String,
		Id:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *database.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.FirstName.String + " " + user.LastName.String
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return next(c)
		}

		u := c.Get("user").(*jwt.Token)

		claims, ok := u.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("failed to cast claims as jwt.MapClaims")
			return next(c)
		}
		expirationFloat := claims["exp"].(float64)
		expiration := int64(expirationFloat)
		if time.Unix(expiration, 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}
				name := claims["name"].(string)
				fmt.Printf(">>>>> name: %v", name)
				idFloat := claims["id"].(float64)
				id := int64(idFloat)
				fmt.Printf(">>>>> id: %v", id)
				if tkn != nil && tkn.Valid {
					_ = GenerateTokensAndSetCookies(&database.User{
						FirstName: sql.NullString{String: strings.Split(name, " ")[0], Valid: true},
						LastName:  sql.NullString{String: strings.Split(name, " ")[1], Valid: true},
						ID:        id,
					}, c)
				}
			}
		}

		return next(c)
	}
}
