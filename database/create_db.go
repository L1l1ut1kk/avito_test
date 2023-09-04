package db

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	DB        *sql.DB
	once      sync.Once
	dbInitErr error
)

func InitDB() error {
	once.Do(func() {
		dbName := "user_segments"
		user := "postgres"
		password := "postgres"

		// Connection string
		connStr := "user=" + user + " password=" + password + " dbname=postgres sslmode=disable"
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Check if the database exists
		var exists bool
		err = conn.QueryRow("SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			// Create the database if it doesn't exist
			_, err = conn.Exec("CREATE DATABASE " + dbName)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Database created successfully")

		// Connect to the user_segments database
		connStr = "user=" + user + " password=" + password + " dbname=" + dbName + " sslmode=disable"
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		// Create tables if they don't exist
		if err := createTablesIfNotExist(); err != nil {
			log.Fatal(err)
		}

		dbInitErr = nil
	})

	return dbInitErr
}

func createTablesIfNotExist() error {
	// Create the "segments" table if it doesn't exist
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS segments (
		id SERIAL PRIMARY KEY,
		slug VARCHAR(255) UNIQUE
	)`)
	if err != nil {
		return err
	}

	// Create the "user_segments" table if it doesn't exist
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS user_segments (
		id SERIAL PRIMARY KEY,
		user_id INT,
		segment_slug VARCHAR(255),
		FOREIGN KEY (segment_slug) REFERENCES segments(slug)
	)`)
	if err != nil {
		return err
	}

	return nil
}
