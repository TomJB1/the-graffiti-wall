package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {

	log.Println("starting")
	r := mux.NewRouter()
	db = connectToDatabase()
	makeTable(db)

	fs := http.FileServer(http.Dir("./website"))

	r.HandleFunc("/", handleAddMessage).Methods("POST")
	r.HandleFunc("/", displayIndex).Methods("GET")

	r.PathPrefix("/").Handler(fs)

	http.ListenAndServe(":1000", r)

}

func handleAddMessage(w http.ResponseWriter, r *http.Request) {

	addMessage(db, r.FormValue("Contents"), r.FormValue("Name"), r.FormValue("Website"), r.FormValue("Email"))
	displayIndex(w, r)

}

func displayIndex(w http.ResponseWriter, r *http.Request) {
	messages := getMessages(db)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, messages)
	if err != nil {
		log.Println(err)
	}
}
