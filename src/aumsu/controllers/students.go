package controllers

import (
	"aumsu/entities"
	models "aumsu/modules"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
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
	r.HandleFunc("/messages/last", getLastMessage).Methods("GET")
	r.HandleFunc("/messages", sendMessage).Methods("POST")
	r.HandleFunc("/messages", getMessages).Methods("GET")
	r.HandleFunc("/messages/{id}", getMessageById).Methods("GET")
	r.HandleFunc("/messages/comment", createComment).Methods("POST")
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

	response, _ := json.Marshal(student)
	w.Write(response)
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

	response, _ := json.Marshal(student)
	w.Write(response)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	var studentModule models.StudentModel
	_, err := studentModule.GetByToken(token)
	if token == "" || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var messageModel models.MessageModel
	messages := messageModel.All()

	response, _ := json.Marshal(messages)
	w.Write(response)
}

func getLastMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	var studentModule models.StudentModel
	student, err := studentModule.GetByToken(token)
	if token == "" || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if student.Status != "user" {
		w.WriteHeader(http.StatusConflict)
		return
	}

	var messageModule models.MessageModel
	message, err := messageModule.GetLast()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(message)
	w.Write(response)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	var studentModule models.StudentModel
	student, err := studentModule.GetByToken(token)
	if token == "" || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" || description == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileName string
	r.ParseMultipartForm(10 << 23)
	file, handler, err := r.FormFile("image")
	if err == nil {
		tempFile, err := ioutil.TempFile("/var/www/images/messages", "image-*-" + filepath.Ext(handler.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tempFile.Write(fileBytes)
		fileName = strings.ReplaceAll(tempFile.Name(), "/var/www/images/messages/", "")
	}

	fmt.Printf("test: " + title)
	message := entities.Message{
		From: student.Id,
		Title: title,
		Description: description,
		Image: fileName,
	}
	var messageModel models.MessageModel
	messageModel.Create(&message)

	pusherClient := pusher.Client{
		AppID:   "966947",
		Key:     "8da04f0e1ecfefbeaecc",
		Secret:  "7d92e3ac99cd7e9e6b3f",
		Cluster: "eu",
	}

	err = pusherClient.Trigger("study-message", "messages", message)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	response, _ := json.Marshal(message)
	w.Write(response)
}

func getMessageById(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]

	var messageModule models.MessageModel
	message, err := messageModule.GetById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(message)
	w.Write(response)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	var comment entities.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var commentModule models.CommentModule
	commentModule.Create(&comment)

	response, _ := json.Marshal(comment)
	w.Write(response)
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

	response, _ := json.Marshal(updatedStudent)
	w.Write(response)
}

func updateAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 23)
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

	response, _ := json.Marshal(updatedStudent)
	w.Write(response)
}
