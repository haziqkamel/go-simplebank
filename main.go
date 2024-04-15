package main

import (
	"database/sql"
	"log"

	"github.com/haziqkamel/simplebank/api"
	db "github.com/haziqkamel/simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	// DefaultPort is the default port the server will listen on
	DefaultPort   = "8080"
	serverAddress = "localhost"
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress + ":" + DefaultPort)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
