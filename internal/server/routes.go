package server

import (
	"net/http"

	"fmt"

	"RemoteMonitor/internal/handlers"
	"RemoteMonitor/static"

	"github.com/gorilla/websocket"
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

	handlers := &handlers.Handler{}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")
	fileServer := http.FileServer(http.FS(static.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", handlers.Login)
	e.GET("/dashboard", handlers.Dashboard)

	/* e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler))) */

	e.GET("/health", s.healthHandler)

	e.GET("/ws", s.websocketHandler)
	e.GET("/wstest", handlers.WsTest)
	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
