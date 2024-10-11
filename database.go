package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Contents string
	Name     string
	Website  string
	Email    string
}

func connectToDatabase() *sql.DB {
	database, _ := sql.Open("sqlite3", "./messages.db")
	err := database.Ping()
	if err != nil {
		log.Println(err)
	}
	log.Println("database may take a minute or two to start up")
	return database
}

func makeTable(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS
	Messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		name TEXT,
		website TEXT,
		email TEXT
    );`)

	log.Println("made table if not exist")

	if err != nil {
		log.Panicln(err)
	}
}

func addMessage(db *sql.DB, content string, name string, website string, email string) {
	log.Println("adding message")
	query := `
    INSERT INTO Messages (content, name, website, email) VALUES (?, ?, ?, ?)
    ;`
	_, err := db.Exec(query, content, name, website, email)

	if err != nil {
		log.Println(err)
	}
}

func getMessages(db *sql.DB) []Message {
	query := `
	SELECT content, name, website, email FROM Messages
	;`
	response, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	var Messages []Message
	for response.Next() {
		var Message Message
		response.Scan(&Message.Contents, &Message.Name, &Message.Website, &Message.Email)
		Messages = append(Messages, Message)
	}
	return Messages
}
