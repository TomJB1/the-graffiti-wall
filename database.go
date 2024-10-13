package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type Message struct {
	Id       int
	Contents string
	Name     string
	Website  string
	Email    string
}

func connectToDatabase() *sql.DB {
	database, _ := sql.Open("sqlite", "./messages.db")
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
		email TEXT,
		upvotes INT,
		time INT
    );`)

	log.Println("made table if not exist")

	if err != nil {
		log.Panicln(err)
	}
}

func removeLowest(db *sql.DB) {
	query := `
    SELECT COUNT(*) FROM Messages
    ;`
	response, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	var rowsStr string
	for response.Next() {
		response.Scan(&rowsStr)
	}

	rows, _ := strconv.Atoi(rowsStr)
	log.Println(rowsStr)
	if rows >= 10 {

		log.Println("removing lowest")
		query = `
		DELETE FROM Messages WHERE id = (SELECT id FROM Messages ORDER BY upvotes ASC, time ASC LIMIT 1)
		;`
		_, err = db.Exec(query)

		if err != nil {
			log.Println(err)
		}
	}

}

func addMessage(db *sql.DB, content string, name string, website string, email string) {
	removeLowest(db)
	log.Println("adding message")
	query := `
    INSERT INTO Messages (content, name, website, email, upvotes, time) VALUES (?, ?, ?, ?, 0, ?)
    ;`
	_, err := db.Exec(query, content, name, website, email, time.Now().Unix())

	if err != nil {
		log.Println(err)
	}
}

func addVote(db *sql.DB, id int) {
	log.Println("adding vote")
	query := `
    UPDATE Messages SET upvotes = upvotes + 1 WHERE id=?
    ;`
	_, err := db.Exec(query, id)

	if err != nil {
		log.Println(err)
	}
}

func getName(db *sql.DB, id int) (string, error) {
	log.Println("selecting name")
	query := `
	SELECT name FROM Messages WHERE id=?
	;`
	response, err := db.Query(query, id)
	if err != nil {
		log.Println(err)
	}
	var name string
	for response.Next() {
		response.Scan(&name)
	}

	return name, err
}

func getMessages(db *sql.DB) []Message {
	query := `
	SELECT id, content, name, website, email FROM Messages ORDER BY upvotes ASC, time ASC
	;`
	response, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	var Messages []Message
	for response.Next() {
		var Message Message
		response.Scan(&Message.Id, &Message.Contents, &Message.Name, &Message.Website, &Message.Email)
		Messages = append(Messages, Message)
	}
	return Messages
}
