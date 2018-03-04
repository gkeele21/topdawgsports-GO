package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbteam"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT TeamID, FullName, Abbreviation, Mascot FROM Team WHERE TeamID > 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var teamid int64
		var displayname, abbrev, mascot database.NullString
		if err := rows.Scan(&teamid, &displayname, &abbrev, &mascot); err != nil {
			log.Fatal(err)
		}

		max := 9;
		if abbrev.String != "" {
			if len(abbrev.String) < max {
				max = len(abbrev.String)
			}
			abbrev.String = abbrev.String[0:max]
		}

		fmt.Printf("TeamID : [%d], Name : [%s], Abbrev : [%s], Mascot : [%s]\n", teamid, displayname, abbrev, mascot)

		team := dbteam.Team{
			TeamID:       teamid,
			Name:         displayname.String,
			Abbreviation: abbrev,
			Mascot:       mascot,
		}

		fmt.Printf("Team : %v\n", team)
		err := dbteam.Insert(&team)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
