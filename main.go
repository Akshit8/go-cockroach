package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var (
	dbUser string
	dbHost string
	dbPort string
)

func init() {
	dbUser = os.Getenv("DB_USER")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	log.Print("dbHost: ", dbHost)
}

func main() {
	config, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s@%s:%s/bank?sslmode=disable", dbUser, dbHost, dbPort))
	if err != nil {
		log.Fatal("error parsing db source: ", err)
	}

	// connect to database
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}
	log.Print("successfully connected to cockroach db")
	defer conn.Close(context.Background())
}
