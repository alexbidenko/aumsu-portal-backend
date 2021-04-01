package main

import (
	"aumsu/controllers"
	"aumsu/dif"
	_ "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello AUMSU!")
}

func main() {
	err := dif.DBError
	if err != nil {
		panic(err)
	}

	dif.Migrate()

	r := mux.NewRouter()

	r.PathPrefix("/files/avatars/").Handler(http.StripPrefix("/files/avatars/", http.FileServer(http.Dir("/var/www/avatars"))))
	r.PathPrefix("/files/messages/images/").Handler(http.StripPrefix("/files/messages/images/", http.FileServer(http.Dir("/var/www/images/messages"))))

	s := r.PathPrefix("/api").Subrouter()
	controllers.InitStudents(s)
	r.HandleFunc("/api", handler)

	fmt.Printf("Server starting")
	log.Fatal(http.ListenAndServe(":8010", r))
}
