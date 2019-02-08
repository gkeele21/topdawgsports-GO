package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbfantasyschedulebase"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ScheduleName, WeekNo, GameNo, Team1, Team2 FROM FSFootballScheduleConfig")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var gameno, team1, team2 int64
		var schedulename, weekno string
		if err := rows.Scan(&schedulename, &weekno, &gameno, &team1, &team2); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Name : [%s], Game : [%s]\n", schedulename, gameno)

		div := dbfantasyschedulebase.FantasyScheduleBase{
			Name:       schedulename,
			WeekName:   weekno,
			GameNumber: gameno,
			HomeTeam:   team1,
			AwayTeam:   team2,
		}

		fmt.Printf("Sched : %#v\n", div)
		err := dbfantasyschedulebase.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
