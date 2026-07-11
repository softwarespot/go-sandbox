package main

import (
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	db, err := sql.Open("duckdb", "data.db?autoinstall_known_extensions=true&autoload_known_extensions=true")
	if err != nil {
		log.Fatal(err)
	}

	// db.Exec(`SET autoinstall_known_extensions=1`)
	// db.Exec(`SET autoload_known_extensions=1`)
	// db.Exec(`CREATE TABLE person (id INTEGER, full_name VARCHAR)`)
	// db.Exec(`INSERT INTO person VALUES (42, 'John Smith')`)

	rows, err := db.Query(`SELECT id, name FROM read_json_auto('in.json')`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		log.Println("got:", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
