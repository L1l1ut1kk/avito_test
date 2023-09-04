package handler_pac

import (
	db "avito/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Segment represents a structure for storing segments
type Segment struct {
	Slug string `json:"slug"`
}

// UserSegments represents a structure for storing user information and their segments
type UserSegments struct {
	UserID   int      `json:"user_id"`
	Segments []string `json:"segments"`
}

// CreateSegment is a handler method for creating a segment.
func CreateSegment(w http.ResponseWriter, r *http.Request) {
	var segment Segment
	if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Call the database method to create a segment
	err := db.CreateSegment(segment.Slug)
	if err != nil {
		log.Printf("Failed to create segment: %v", err)
		http.Error(w, "Failed to create segment", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{"Segment created successfully"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteSegment is a handler method for deleting a segment
func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	// Call the database method to delete a segment
	err := db.DeleteSegment(slug)
	if err != nil {
		log.Printf("Failed to delete segment: %v", err)
		http.Error(w, "Failed to delete segment", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{"Segment deleted successfully"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddUserToSegment is a handler method for adding a user to a segment
func AddUserToSegment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDStr := params["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var userSegmentsData UserSegments
	if err := json.NewDecoder(r.Body).Decode(&userSegmentsData); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Get the segments already associated with the user
	currentSegments, err := db.GetActiveUserSegments(userID)
	if err != nil {
		log.Printf("Failed to get user's active segments: %v", err)
		http.Error(w, "Failed to get user's active segments", http.StatusInternalServerError)
		return
	}
	for _, segment := range userSegmentsData.Segments {
		err := db.RegisterUserSegmentOperation(userID, segment, "add")
		if err != nil {
			log.Printf("Failed to register user segment operation: %v", err)
			http.Error(w, "Failed to register user segment operation", http.StatusInternalServerError)
			return
		}
	}

	// Create a map to track existing segments for faster lookup
	existingSegments := make(map[string]bool)
	for _, segment := range currentSegments {
		existingSegments[segment] = true
	}

	// Create a map to track added segments for deduplication
	addedSegments := make(map[string]bool)

	// Add the user to segments in the database if they are not already associated
	for _, segment := range userSegmentsData.Segments {
		if _, exists := existingSegments[segment]; !exists {
			if _, alreadyAdded := addedSegments[segment]; !alreadyAdded {
				err := db.AddUserToSegment(userID, segment)
				if err != nil {
					log.Printf("Failed to add user to segment: %v", err)
					http.Error(w, "Failed to add user to segment", http.StatusInternalServerError)
					return
				}
				addedSegments[segment] = true
			}
		}
	}

	// Construct the response message
	var responseMessage string
	if len(addedSegments) > 0 {
		responseMessage = "User added to segments successfully: " + strings.Join(getKeys(addedSegments), ", ")
	} else {
		responseMessage = "No changes were made. User is already associated with all specified segments."
	}

	response := struct {
		Message string `json:"message"`
	}{responseMessage}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusCreated)
}

// Take key function
func getKeys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// GetActiveSegments is a handler method for getting active segments of a user
func GetActiveSegments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDStr := params["user_id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the database method to get active segments of a user
	activeSegments, err := db.GetActiveUserSegments(userID)
	if err != nil {
		log.Printf("Failed to get active segments: %v", err)
		http.Error(w, "Failed to get active segments", http.StatusInternalServerError)
		return
	}

	// Sending active segments in response
	json.NewEncoder(w).Encode(UserSegments{UserID: userID, Segments: activeSegments})
}

// DeleteUserSegments is a handler method for deleting segments from a user
func DeleteUserSegments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDStr := params["user_id"]
	//segment := params["segment_slug"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var segmentsToDelete UserSegments
	if err := json.NewDecoder(r.Body).Decode(&segmentsToDelete); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Обновляем таблицу user_segment_history для фиксации удаляемых сегментов
	for _, deletedSegment := range segmentsToDelete.Segments {
		err = db.AddUserSegmentHistory(userID, deletedSegment, "delit")
		if err != nil {
			log.Printf("Failed to add user segment history: %v", err)
			http.Error(w, "Failed to add user segment history", http.StatusInternalServerError)
			return
		}
	}

	// Затем удалите сегменты из базы данных, как вы это делали ранее
	err = deleteSegmentsForUser(userID, segmentsToDelete.Segments)
	if err != nil {
		log.Printf("Failed to delete segments for user: %v", err)
		http.Error(w, "Failed to delete segments for user", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{"Segments deleted successfully"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteSegmentsForUser(userID int, segmentsToDelete []string) error {
	for _, segment := range segmentsToDelete {
		err := db.RemoveUserFromSegment(userID, segment)
		if err != nil {
			return err
		}
	}
	return nil
}
