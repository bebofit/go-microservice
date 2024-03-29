package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {
	log.Println("Starting authentication service")
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres!")
	}
	app := Config{
		Client: &http.Client{},
	}
	app.setupRepo(conn)
	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, error := openDB(dsn)
		if error != nil {
			log.Println("Postgres not ready yet ...")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}
		if counts > 10 {
			log.Println(error)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepo(conn)
	app.Repo = db
}
