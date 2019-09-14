package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbconference"
)

var db *sql.DB

func main() {
	// grab all conferences from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(topdawg.circlepix.com:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ConferenceID, DisplayName, FullName FROM Conference")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var conferenceid int64
		var name, abbrev database.NullString
		if err := rows.Scan(&conferenceid, &abbrev, &name); err != nil {
			log.Fatal(err)
		}

		max := 9
		if abbrev.String != "" {
			if len(abbrev.String) < max {
				max = len(abbrev.String)
			}
			abbrev.String = abbrev.String[0:max]
		}

		fmt.Printf("ConferenceID : [%d], Name : [%s], Abbrev : [%s]\n", conferenceid, name, abbrev)

		conf := dbconference.Conference{
			ConferenceID: conferenceid,
			Name:         name.String,
			Abbreviation: abbrev,
		}

		fmt.Printf("Conf : %v\n", conf)
		err := dbconference.Insert(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
