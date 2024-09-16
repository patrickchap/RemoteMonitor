package server

import (
	"errors"
	"net/http"

	"fmt"

	hd "RemoteMonitor/internal/handlers"
	"RemoteMonitor/static"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

func (s *Server) websocketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func (s *Server) RegisterRoutes() http.Handler {

	handlers := &hd.Handler{
		Store: s.Store,
	}
	handlers.NewSession(hd.CookieOpts{
		Name:   "auth-session",
		Secret: "super-secret-key",
		MaxAge: 86400 * 7,
	}, func(c echo.Context) error {
		return c.Redirect(http.StatusFound,
			c.Echo().Reverse("/"))
	})
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")
	fileServer := http.FileServer(http.FS(static.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", handlers.Login)

	/* e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler))) */

	e.GET("/ws", s.websocketHandler)
	e.GET("/wstest", handlers.WsTest)

	e.POST("/login", handlers.PostLogin)
	e.POST("/logout", handlers.Logout)
	adminGroup := e.Group("/admin")
	adminGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(hd.GetJWTSecret()),
		TokenLookup: "cookie:access-token",
		ErrorHandler: func(c echo.Context, err error) error {
			_, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
			if !ok {
				fmt.Println("Redirecting to login")
				target := c.Request().URL.Path
				return c.Redirect(http.StatusFound, fmt.Sprintf("/?target=%s", target))
			}
			return nil
		},
	}))
	adminGroup.Use(hd.TokenRefresherMiddleware)
	adminGroup.GET("/", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}
		fmt.Println(claims.GetExpirationTime())
		name := claims["name"].(string)
		idFloat := claims["id"].(float64)
		id := int64(idFloat)
		fmt.Println(name, id)
		return c.JSON(http.StatusOK, claims)
	})

	adminGroup.GET("/dashboard", handlers.Dashboard)
	adminGroup.GET("/hosts", handlers.Hosts)
	adminGroup.GET("/host/edit/:id", handlers.HostEdit)
	adminGroup.GET("/host/edit/form/:id", handlers.GetEditHostDetails)

	adminGroup.PUT("/host/edit/form", handlers.PutEditHostDetails)

	adminGroup.GET("/host/create", handlers.HostCreateForm)
	adminGroup.POST("/host/create", handlers.HostCreate)

	adminGroup.GET("/host/edit/hostservice/:id", handlers.GetHostServices)
	adminGroup.POST("/hostservice/create", handlers.PostHostService)
	adminGroup.GET("/hostservice/edit/:id", handlers.EditServiceRow)
	adminGroup.GET("/hostservice/edit/row/:id", handlers.GetServiceRow)
	adminGroup.DELETE("/hostservice/delete/:id", handlers.DeleteServiceRow)
	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
