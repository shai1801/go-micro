package main

import (
	"authentication/data"
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

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the authentication service")

	//TODO connect to DB

	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to DB")
	}

	//set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	//init the server

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start the server
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

// opening the DB connection.
func openDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {

		return nil, err
	}

	return db, nil
}

// It might happen that authentication service starts befor the PostgresDB
func connectToDB() *sql.DB {

	dsn := os.Getenv("DSN")
	for {
		db, err := openDB(dsn)
		if err != nil {

			fmt.Println("Data base is not up.")
			count++

		} else {
			log.Println("db is connected")
			return db
		}
		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Retrying to connect after 2 seconds")
		time.Sleep(2 * time.Second)

	}

}
