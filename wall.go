package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Pagedata struct {
	Messages         []Message
	Vote             string
	SubmittedMessage bool
}

var db *sql.DB

func main() {

	log.Println("starting")
	r := mux.NewRouter()
	db = connectToDatabase()
	makeTable(db)

	fs := http.FileServer(http.Dir("./website"))

	r.HandleFunc("/add", handleAddMessage).Methods("POST")
	r.HandleFunc("/vote", handleVote).Methods("POST")
	r.HandleFunc("/", displayIndex).Methods("GET")

	r.PathPrefix("/").Handler(fs)

	http.ListenAndServe(":1000", r)

}

func handleAddMessage(w http.ResponseWriter, r *http.Request) {
	_, cookie_err := r.Cookie("message-submitted") // cookie_err should NOT be nil to add message
	if cookie_err != nil {
		addMessage(db, r.FormValue("Contents"), r.FormValue("Name"), r.FormValue("Website"), r.FormValue("Email"))
		messageCookie := &http.Cookie{Name: "message-submitted", Value: "true", HttpOnly: false, Expires: time.Now().Add(24 * time.Hour)}
		http.SetCookie(w, messageCookie)

	} else {
		println("already submitted message")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handleVote(w http.ResponseWriter, r *http.Request) {
	vote, err := strconv.Atoi(r.FormValue("Vote"))
	_, cookie_err := r.Cookie("vote") // cookie_err should NOT be nil to vote
	if err == nil && cookie_err != nil {
		addVote(db, vote)
		voteCookie := &http.Cookie{Name: "vote", Value: r.FormValue("Vote"), HttpOnly: false, Expires: time.Now().Add(24 * time.Hour)}
		http.SetCookie(w, voteCookie)

	} else {
		println("invalid vote")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func displayIndex(w http.ResponseWriter, r *http.Request) {
	var pagedata Pagedata

	voteCookie, voteCookie_err := r.Cookie("vote")
	messages := getMessages(db)
	if voteCookie_err == nil {
		vote, convert_err := strconv.Atoi(voteCookie.Value)
		if convert_err != nil {
			log.Println(convert_err)
		}

		name, name_err := getName(db, vote)
		if name_err != nil {
			log.Println(name_err)
		}

		pagedata.Vote = name
	}

	_, messageCookie_err := r.Cookie("message-submitted")
	if messageCookie_err != nil { // there is NOT a cookie
		pagedata.SubmittedMessage = false
	} else {
		pagedata.SubmittedMessage = true
	}

	pagedata.Messages = messages
	tmpl := template.Must(template.ParseFiles("website/templates/index.html"))
	err := tmpl.Execute(w, pagedata)
	if err != nil {
		log.Println(err)
	}
}
