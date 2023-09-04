package main

import (
	db "avito/database"
	handler_pac "avito/handlers"
	"log"
	"net/http"

	_ "avito/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

// @title           Avito test task
// @version         1.0
// @description     Task for internship
// @contact.name   Людмила Мишакова l1l1ut1kk@mail.ru
// @license.name  Ubuntu 22.04
// @host      localhost:8080

// @BasePath
func main() {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("docs/*any"),
	))
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
