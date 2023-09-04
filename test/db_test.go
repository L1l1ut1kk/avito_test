package test

import (
	db "avito/database"
	"testing"
)

func TestInitDB(t *testing.T) {
	// Calling InitDB and checking for errors
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Checking that db is not equal to nil
	if db.GetDB() == nil {
		t.Fatal("Expected db to be initialized, but it is nil")
	}
}

func TestGetDB(t *testing.T) {
	// Ensure that the database is initialized before the test
	db.InitDB()

	// Get an active database connection using GetDB
	database := db.GetDB()

	// Checking that the returned connection is not equal to nil
	if database == nil {
		t.Fatal("Expected an active database connection, but got nil")
	}

	// Try executing an SQL query to ensure that the connection is working
	_, err := database.Exec("SELECT 1")
	if err != nil {
		t.Fatalf("Expected the database connection to work, but got an error: %v", err)
	}
}
