package controllers

import (
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	"aumsu.portal.backend/utils"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type Authorization struct {
	Login    string
	Password string
}

func InitStudents(r *mux.Router) {
	r.HandleFunc("/login", authorization).Methods("POST")
	r.HandleFunc("/registration", registration).Methods("POST")
	r.HandleFunc("/user", updateStudent).Methods("PUT")
	r.HandleFunc("/user/avatar", updateAvatar).Methods("PUT")
}

func authorization(w http.ResponseWriter, r *http.Request) {
	var data Authorization
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var studentModule models.StudentModel
	student, err := studentModule.Authorization(data.Login, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJsonResponse(w, student)
}

func registration(w http.ResponseWriter, r *http.Request) {
	var data entities.Student
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var studentModule models.StudentModel

	existed := studentModule.CheckUnique(&data)
	if existed {
		http.Error(w, "Entity is existed", http.StatusBadRequest)
		return
	}

	student := entities.Student{
		Login: data.Login,
		Token: data.Login,
		Password: data.Password,
		FirstName: data.FirstName,
		LastName: data.LastName,
		Patronymic: data.Patronymic,
		Status: "user",
		Avatar: "",
	}
	studentModule.Create(&student)

	utils.WriteJsonResponse(w, student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	var student entities.Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var studentModule models.StudentModel
	updatedStudent, _ := studentModule.GetByToken(r.Header.Get("Authorization"))
	updatedStudent.FirstName = student.FirstName
	updatedStudent.LastName = student.LastName
	updatedStudent.Patronymic = student.Patronymic
	studentModule.Update(updatedStudent.Id, &updatedStudent)

	utils.WriteJsonResponse(w, updatedStudent)
}

func updateAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 23)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileName string
	tempFile, err := ioutil.TempFile("/var/www/images/avatars", "avatar-*-" + filepath.Ext(handler.Filename))
	if err != nil {
		http.Error(w, "Create temporary file: " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Read file: " + err.Error(), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)
	fileName = strings.ReplaceAll(tempFile.Name(), "/var/www/images/avatars/", "")

	var studentModule models.StudentModel
	updatedStudent, _ := studentModule.GetByToken(r.Header.Get("Authorization"))
	updatedStudent.Avatar = fileName
	studentModule.Update(updatedStudent.Id, &updatedStudent)

	utils.WriteJsonResponse(w, updatedStudent)
}
