package main

import (
	db "avito/database"
	handler_pac "avito/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db.InitDB()

	// HTTP request
	r.HandleFunc("/segments", handler_pac.CreateSegment).Methods("POST")
	r.HandleFunc("/segments/{slug}", handler_pac.DeleteSegment).Methods("DELETE")
	r.HandleFunc("/users/{user_id}/segments", handler_pac.AddUserToSegment).Methods("POST")
	r.HandleFunc("/users/{user_id}/segments", handler_pac.GetActiveSegments).Methods("GET")
	r.HandleFunc("/users/{user_id}/segments", handler_pac.DeleteUserSegments).Methods("DELETE")

	r.HandleFunc("/generate-report/{year}/{month}", handler_pac.GenerateReport).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
