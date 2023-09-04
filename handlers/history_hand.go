package handler_pac

import (
	// Импорт необходимых библиотек и пакетов
	db "avito/database"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GenerateReport create file in CSV in select period
func GenerateReport(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	yearStr := params["year"]
	monthStr := params["month"]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	// Select data
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	operations, err := db.GetOperationsByDateRange(startDate, endDate)
	if err != nil {
		log.Printf("Failed to retrieve operations: %v", err)
		http.Error(w, "Failed to retrieve operations", http.StatusInternalServerError)
		return
	}

	// Create CSV-file
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=report.csv")

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()
	header := []string{"UserID", "Segment", "Operation", "Timestamp"}
	csvWriter.Write(header)

	// Write operation in CSV
	for _, op := range operations {
		row := []string{
			strconv.Itoa(op.UserID),
			op.Segment,
			op.Operation,
			op.Timestamp.Format("2006-01-02 15:04:05"),
		}
		csvWriter.Write(row)
	}
}
