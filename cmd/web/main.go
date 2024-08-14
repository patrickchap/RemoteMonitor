package main

import (
	database "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/server"
	"database/sql"
	"fmt"
	"log"
	"os"

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

	server := server.NewServer(store)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
