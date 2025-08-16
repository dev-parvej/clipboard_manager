package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS clipboard (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		image_path TEXT,
		timestamp DATETIME
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func StoreTextClipboard(txt string) {
	_, err := db.Exec("INSERT INTO clipboard(content, timestamp) VALUES(?, ?)", txt, time.Now())
	if err != nil {
		log.Println("Error storing text:", err)
	}
}

func StoreImageClipboard(imgPath string) {
	_, err := db.Exec("INSERT INTO clipboard(image_path, timestamp) VALUES(?, ?)", imgPath, time.Now())
	if err != nil {
		log.Println("Error storing image:", err)
	}
}

func CleanupOldEntries() {
	threshold := time.Now().AddDate(0, 0, -30)
	_, err := db.Exec("DELETE FROM clipboard WHERE timestamp < ?", threshold)
	if err != nil {
		log.Println("Error cleaning old entries:", err)
	}
}
