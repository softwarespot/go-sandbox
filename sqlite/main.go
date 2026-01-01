package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// URL: https://github.com/philippta/sqlitebench
// URL: https://github.com/gwenn/gosqlite
// URL. https://github.com/benbjohnson/wtf
// URL: https://github.com/xo/usql
// URL: https://github.com/golang-migrate/migrate

func main() {
	db, err := sql.Open("sqlite", "file:data.db?_pragma=journal_mode=WAL&_pragma=synchronous=normal&cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	concurrency := runtime.NumCPU() / 2
	db.SetMaxOpenConns(concurrency)

	create(db)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go read(db, i, 1_000_000)
	}
	wg.Wait()
}

func create(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS "testing" (
			"id"	INTEGER,
			"value"	INTEGER,
			PRIMARY KEY("id" AUTOINCREMENT)
		);
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func insert(db *sql.DB, id int) {
	t0 := time.Now()
	query := `
		INSERT OR REPLACE INTO testing
		(id, value)
		VALUES (?, ?)
	`

	if _, err := db.Exec(query, id, id); err != nil {
		log.Fatal(err, "inserting ID:", strconv.Itoa(id))
	}
	log.Printf("INSERT %d, took %v\n", id, time.Since(t0))
}

func read(db *sql.DB, clientID, maxID int) {
	for {
		id := rand.Intn(maxID)

		// There is no ID of 0 in SQLite
		if id == 0 {
			id = 1
		}

		query := `
			SELECT value
			FROM testing
			WHERE id = ?
		`
		var value int
		t0 := time.Now()
		err := db.QueryRow(query, id).Scan(&value)
		msg := fmt.Sprintf("client ID: %d, ID: %d, value: %d, took: %v", clientID, id, value, time.Since(t0))

		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows.", msg)
			insert(db, id)
			continue
		}

		if err != nil {
			log.Fatal(err, "get ID:", strconv.Itoa(id))
		}

		log.Println(msg)
		time.Sleep(16 * time.Millisecond)
	}
}
