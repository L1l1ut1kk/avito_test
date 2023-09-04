package db

import (
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Operation struct { // Define the Operation struct
	UserID    int       `json:"user_id"`
	Segment   string    `json:"segment"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

// Add segment to history
func AddUserSegmentHistory(userID int, segmentSlug string, operation string) error {
	db := GetDB()
	_, err := db.Exec("INSERT INTO user_segment_history (user_id, segment_slug, operation) VALUES ($1, $2, $3)", userID, segmentSlug, operation)
	return err
}

func RegisterUserSegmentOperation(userID int, segmentSlug string, operation string) error {
	// Function for add params in table user_segment_history
	err := AddUserSegmentHistory(userID, segmentSlug, operation)
	if err != nil {
		log.Printf("Failed to register user segment operation: %v", err)
		return err
	}
	return nil
}

func GetOperationsByDateRange(startDate, endDate time.Time) ([]Operation, error) {
	db := GetDB() // Replace GetDB with your actual database connection setup

	rows, err := db.Query(`
        SELECT user_id, segment_slug, operation, timestamp
        FROM user_segment_history
        WHERE timestamp >= $1 AND timestamp <= $2
    `, startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []Operation
	for rows.Next() {
		var op Operation
		if err := rows.Scan(&op.UserID, &op.Segment, &op.Operation, &op.Timestamp); err != nil {
			return nil, err
		}
		operations = append(operations, op)
	}

	return operations, nil
}
