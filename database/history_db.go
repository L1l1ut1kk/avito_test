package db

import (
	"log"
	"time"

	_ "github.com/lib/pq"
)

type UserSegmentHistory struct {
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
	// Вызовите функцию для добавления записи в таблицу user_segment_history
	err := AddUserSegmentHistory(userID, segmentSlug, operation)
	if err != nil {
		log.Printf("Failed to register user segment operation: %v", err)
		return err
	}
	return nil
}

// Back users history
func GetActiveUserSegmentsWithinDateRange(userID int, startDate, endDate time.Time) ([]string, error) {
	db := GetDB()

	rows, err := db.Query(`
        SELECT s.slug
        FROM segments s
        JOIN user_segments us ON s.slug = us.segment_slug
        WHERE us.user_id = $1
        AND us.timestamp >= $2
        AND us.timestamp <= $3
    `, userID, startDate, endDate)

	if err != nil {
		log.Printf("Failed to fetch active user segments within date range: %v", err)
		return nil, err
	}
	defer rows.Close()

	var activeSegments []string
	for rows.Next() {
		var segmentSlug string
		if err := rows.Scan(&segmentSlug); err != nil {
			log.Printf("Error scanning segment slug: %v", err)
			return nil, err
		}
		activeSegments = append(activeSegments, segmentSlug)
	}

	return activeSegments, nil
}
