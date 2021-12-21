package controllers

import (
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	"aumsu.portal.backend/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Authorization struct {
	Login    string
	Password string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func InitStudents(r *mux.Router) {
	r.HandleFunc("/version/{number}", getVersion).Methods("GET")
	r.HandleFunc("/login", authorization).Methods("POST")
	r.HandleFunc("/registration", registration).Methods("POST")
	r.HandleFunc("/user", getStudent).Methods("GET")
	r.HandleFunc("/schedule/{id}", getSchedule).Methods("GET")
	r.HandleFunc("/user", updateStudent).Methods("PUT")
	r.HandleFunc("/user/avatar", updateAvatar).Methods("PUT")
	r.HandleFunc("/user/password", updatePassword).Methods("PUT")
	r.HandleFunc("/study-groups", getStudyGroups).Methods("GET")
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.ParseInt(mux.Vars(r)["number"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJsonResponse(w, x >= 18)
}

func authorization(w http.ResponseWriter, r *http.Request) {
	var data Authorization
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var studentModule models.StudentModel
	student, err := studentModule.Authorization(data.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// TODO: Выпилить
	if strings.HasPrefix(data.Password, "$2a$14$") {
		if student.Password != data.Password {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Println(student.Password, "; ", data.Password, "; ", data.Login, "; ", student, "; ", data)
		err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(data.Password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	utils.WriteJsonResponse(w, student)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	var studentModule models.StudentModel
	updatedStudent, err := studentModule.GetByToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJsonResponse(w, updatedStudent)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var studyGroupModel models.StudyGroupModel
	schedule, err := studyGroupModel.GetSchedule(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJsonResponse(w, schedule)
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

	bytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	student := entities.Student{
		Login:      data.Login,
		Token:      generateString(40),
		Password:   string(bytes),
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		Patronymic: data.Patronymic,
		Status:     "user",
		Avatar:     "",
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

	fmt.Println(student.StudyGroupId)
	var studentModule models.StudentModel
	updatedStudent, _ := studentModule.GetByToken(r.Header.Get("Authorization"))
	updatedStudent.Login = student.Login
	updatedStudent.FirstName = student.FirstName
	updatedStudent.LastName = student.LastName
	updatedStudent.Patronymic = student.Patronymic
	updatedStudent.StudyGroupId = student.StudyGroupId
	studentModule.Update(updatedStudent.Id, &updatedStudent)

	utils.WriteJsonResponse(w, updatedStudent)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	var password = r.FormValue("password")
	var newPassword = r.FormValue("new_password")

	if len(newPassword) < 8 {
		http.Error(w, "Password is short", http.StatusBadRequest)
		return
	}

	var studentModule models.StudentModel
	updatedStudent, _ := studentModule.GetByToken(r.Header.Get("Authorization"))

	err := bcrypt.CompareHashAndPassword([]byte(updatedStudent.Password), []byte(password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedStudent.Password = string(bytes)
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
	tempFile, err := ioutil.TempFile("/var/www/images/avatars", "avatar-*"+filepath.Ext(handler.Filename))
	if err != nil {
		http.Error(w, "Create temporary file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Read file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)
	fileName = strings.ReplaceAll(tempFile.Name(), "/var/www/images/avatars/", "")

	var studentModule models.StudentModel
	updatedStudent, _ := studentModule.GetByToken(r.Header.Get("Authorization"))

	if updatedStudent.Avatar != "" {
		os.Remove("/var/www/images/avatars/" + updatedStudent.Avatar)
	}

	updatedStudent.Avatar = fileName
	studentModule.Update(updatedStudent.Id, &updatedStudent)

	utils.WriteJsonResponse(w, updatedStudent)
}

func getStudyGroups(w http.ResponseWriter, _ *http.Request) {
	var studyGroupModule models.StudyGroupModel
	studyGroups := studyGroupModule.All()

	utils.WriteJsonResponse(w, studyGroups)
}
