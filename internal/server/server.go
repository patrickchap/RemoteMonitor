package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"RemoteMonitor/config"
	database "RemoteMonitor/internal/database/sqlc"
)

type Server struct {
	port      int
	Store     database.Store
	AppConfig *config.AppConfig
}

func NewServer(store database.Store, appConfig *config.AppConfig) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:      port,
		Store:     store,
		AppConfig: appConfig,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
