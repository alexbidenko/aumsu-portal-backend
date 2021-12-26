package main

import (
	"aumsu.portal.backend/controllers"
	"aumsu.portal.backend/dif"
	_ "database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	err := dif.DBError
	if err != nil {
		panic(err)
	}

	dif.Migrate()

	r := mux.NewRouter()

	r.PathPrefix("/files/avatars/").Handler(http.StripPrefix("/files/avatars/", http.FileServer(http.Dir("/var/www/images/avatars"))))
	r.PathPrefix("/files/messages/images/").Handler(http.StripPrefix("/files/messages/images/", http.FileServer(http.Dir("/var/www/images/messages"))))

	s := r.PathPrefix("/api").Subrouter()
	controllers.InitStudents(s)
	controllers.InitMessages(s)

	allowOrigins := handlers.AllowedOrigins([]string{"*"})
	allowHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	allowMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	fmt.Println("Server starting")
	log.Fatal(http.ListenAndServe("0.0.0.0:8010", handlers.CORS(allowOrigins, allowHeaders, allowMethods)(r)))
}
