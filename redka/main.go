// A basic example of using Redka
// with github.com/mattn/go-sqlite3 driver.
package main

import (
	"database/sql"
	"log"

	"github.com/nalgeon/redka"
	driver "modernc.org/sqlite"
)

func main() {
	// modernc.org/sqlite uses a different driver name ("sqlite"),
	// while Redka expects "sqlite3". So we have to re-register it as "sqlite3"
	sql.Register("sqlite3", &driver.Driver{})

	db, err := redka.Open("data.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println(db.Hash().Set("2024-12-01", "day1", "1000"))
	log.Println(db.Hash().Set("2024-12-02", "day2", "2000"))
	log.Println(db.Hash().Set("2024-12-03", "day3", "3000"))
	log.Println(db.Hash().Get("2024-12-02", "day2"))
	log.Println(db.Hash().Get("2024-12-04", "day4"))

	log.Println()

	log.Println(db.List().PushFront("2024-12-04", "day4"))
	log.Println(db.List().Get("2024-12-04", 0))
	log.Println(db.List().Set("2024-12-04", -2, "day4"))

	// db.Key().DeleteAll()
}
