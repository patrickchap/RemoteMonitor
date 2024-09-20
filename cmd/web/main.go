package main

import (
	"RemoteMonitor/config"
	database "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/server"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var dburl = os.Getenv("DB_URL")

func main() {

	sqlDb, err := sql.Open("libsql", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	store := database.NewStore(sqlDb)

	appConfig := &config.AppConfig{}
	timezone, err := time.LoadLocation("Local")
	if err != nil {
		panic(fmt.Sprintf("cannot load location: %s", err))
	}

	scheduler := cron.New(cron.WithLocation(timezone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	appConfig.Schedual = scheduler

	server := server.NewServer(store, appConfig)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
