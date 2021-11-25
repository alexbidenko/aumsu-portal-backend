package main

import (
	"aumsu.portal.backend/controllers"
	"aumsu.portal.backend/dif"
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	_ "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func main() {
	err := dif.DBError
	if err != nil {
		panic(err)
	}

	dif.Migrate()

	var student entities.Student
	dif.DB.Model(&entities.Student{}).Where("login = ?", "alexbidenko").First(&student)
	if student.Password == "12345678" {
		var students []entities.Student
		dif.DB.Model(&entities.Student{}).Find(&students)

		for _, item := range students {
			bytes, _ := bcrypt.GenerateFromPassword([]byte(item.Password), 14)
			item.Password = string(bytes)

			var studentModule models.StudentModel
			studentModule.Update(item.Id, &item)
		}
	}

	r := mux.NewRouter()

	r.PathPrefix("/files/avatars/").Handler(http.StripPrefix("/files/avatars/", http.FileServer(http.Dir("/var/www/images/avatars"))))
	r.PathPrefix("/files/messages/images/").Handler(http.StripPrefix("/files/messages/images/", http.FileServer(http.Dir("/var/www/images/messages"))))

	s := r.PathPrefix("/api").Subrouter()
	controllers.InitStudents(s)
	controllers.InitMessages(s)

	fmt.Printf("Server starting")
	log.Fatal(http.ListenAndServe(":8010", r))
}
