package controllers

import (
	"aumsu.portal.backend/entities"
	models "aumsu.portal.backend/modules"
	"aumsu.portal.backend/utils"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
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
		http.Error(w, "Title and description are required", http.StatusBadRequest)
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
		tempFile, err := ioutil.TempFile("/var/www/images/messages", "image-*-"+filepath.Ext(handler.Filename))
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

	message := entities.Message{
		From:        student.ID,
		Title:       title,
		Description: description,
		Image:       fileName,
	}
	var messageModel models.MessageModel
	messageModel.Create(&message)

	notificationMessage := &messaging.Message{
		Data: map[string]string{
			"sender_id": strconv.Itoa(int(message.From)),
		},
		Notification: &messaging.Notification{
			Title: message.Title,
			Body:  message.Description,
		},
		Topic: "messages",
	}

	opt := option.WithCredentialsFile("aumsu-portal-firebase-adminsdk-5sajn-e6d3adfd5a.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.Send(context.Background(), notificationMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
