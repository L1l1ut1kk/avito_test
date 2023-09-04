package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// GetDB returns the active database connection
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database is not initialized. Call InitDB() before use.")
	}
	return db
}
func RemoveUserFromSegment(userID int, segmentSlug string) error {
	db := GetDB()
	_, err := db.Exec("DELETE FROM user_segments WHERE user_id = $1 AND segment_slug = $2", userID, segmentSlug)
	return err
}

// CreateSegment inserts a new segment into the database.
func CreateSegment(slug string) error {
	db := GetDB()
	_, err := db.Exec("INSERT INTO segments (slug) VALUES ($1)", slug)
	return err
}

// DeleteSegment deletes a segment from the database based on its slug
func DeleteSegment(slug string) error {
	db := GetDB()
	_, err := db.Exec("DELETE FROM segments WHERE slug = $1", slug)
	return err
}

// AddUserToSegment adds a user to a segment in the database
func AddUserToSegment(userID int, segmentSlug string) error {
	db := GetDB()
	_, err := db.Exec("INSERT INTO user_segments (user_id, segment_slug) VALUES ($1, $2)", userID, segmentSlug)
	return err
}

// GetActiveUserSegments returns the active segments for a specified user
func GetActiveUserSegments(userID int) ([]string, error) {
	db := GetDB()
	rows, err := db.Query("SELECT s.slug FROM segments s JOIN user_segments us ON s.slug = us.segment_slug WHERE us.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeSegments []string
	for rows.Next() {
		var segmentSlug string
		if err := rows.Scan(&segmentSlug); err != nil {
			return nil, err
		}
		activeSegments = append(activeSegments, segmentSlug)
	}
	return activeSegments, nil
}
