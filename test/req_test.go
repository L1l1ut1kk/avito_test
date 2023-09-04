package test

import (
	db "avito/database"
	handler_pac "avito/handlers"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

func TestCreateSegment(t *testing.T) {
	db.InitDB()
	// Create a test server and an HTTP request
	server := httptest.NewServer(http.HandlerFunc(handler_pac.CreateSegment))
	defer server.Close()

	segment := handler_pac.Segment{Slug: "test12-segment"}
	jsonSegment, _ := json.Marshal(segment)

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(jsonSegment))
	if err != nil {
		t.Fatal(err)
	}

	// Send an HTTP request to the server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestDeleteSegment(t *testing.T) {
	db.InitDB()
	// Create a test server and an HTTP request with {slug} substitution
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := map[string]string{"slug": "test3-segment"}
		r = mux.SetURLVars(r, vars)
		handler_pac.DeleteSegment(w, r)
	}))
	defer server.Close()

	req, err := http.NewRequest("DELETE", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send an HTTP request to the server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAdd_GetActiveSegments(t *testing.T) {
	// Initialize the test database
	db.InitDB()

	// Create an HTTP request to add a user to a segment
	userID := 1
	segmentSlug := "segment-1"
	err := db.AddUserToSegment(userID, segmentSlug)

	if err != nil {
		t.Fatalf("Failed to add user to segment: %v", err)
	}
	defer cleanUpTestData(userID, segmentSlug)

	// Create an HTTP request to get the active segments of a user
	req, err := http.NewRequest("GET", "/users/1/segments", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Create a test HTTP ResponseWriter
	recorder := httptest.NewRecorder()

	// Call the GetActiveSegments handler
	handler := http.HandlerFunc(handler_pac.GetActiveSegments)
	handler.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	// Parse the JSON response
	var userSegmentsData handler_pac.UserSegments
	if err := json.Unmarshal(recorder.Body.Bytes(), &userSegmentsData); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check the active segments of the user
	if len(userSegmentsData.Segments) != 1 || userSegmentsData.Segments[0] != segmentSlug {
		t.Errorf("Expected active segments to contain [%s], but got %v", segmentSlug, userSegmentsData.Segments)
	}
}

func cleanUpTestData(userID int, segmentSlug string) {
	// Delete data in the database
	if err := db.DeleteSegment(segmentSlug); err != nil {
		log.Printf("Failed to delete segment: %v", err)
	}
}
