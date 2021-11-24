package controllers

import (
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	"aumsu.portal.backend/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func InitMessages(r *mux.Router) {
	r.HandleFunc("/messages/last", getLastMessage).Methods("GET")
	r.HandleFunc("/messages", sendMessage).Methods("POST")
	r.HandleFunc("/messages", getMessages).Methods("GET")
	r.HandleFunc("/messages/{id}", getMessageById).Methods("GET")
	r.HandleFunc("/messages/comment", createComment).Methods("POST")
	r.HandleFunc("/messages/comment/{id}", deleteComment).Methods("DELETE")
	r.HandleFunc("/messages/comment/{id}", updateComment).Methods("PUT")
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

	utils.WriteJsonResponse(w, messages)
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

	utils.WriteJsonResponse(w, message)
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
	err = r.ParseMultipartForm(10 << 23)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	utils.WriteJsonResponse(w, message)
}

func getMessageById(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]

	var messageModule models.MessageModel
	message, err := messageModule.GetById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utils.WriteJsonResponse(w, message)
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

	utils.WriteJsonResponse(w, comment)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	var commentModule models.CommentModule
	commentModule.Delete(mux.Vars(r)["id"])

	utils.WriteJsonResponse(w, true)
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	var data entities.Comment
	utils.ParseRequestBody(w, r, &data)

	var commentModule models.CommentModule
	commentModule.Update(mux.Vars(r)["id"], &data)

	utils.WriteJsonResponse(w, data)
}
