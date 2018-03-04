package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database/dbfantasyplayervalue"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT pv.PlayerID, pv.FSSeasonWeekID, pv.Value FROM FSPlayerValue pv INNER JOIN Player p on p.PlayerID = pv.PlayerID " +
		"INNER JOIN FSSeasonWeek fsw on fsw.FSSeasonWeekID = pv.FSSeasonWeekID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var playerid, fsseasonweekid int64
		var value float64
		if err := rows.Scan(&playerid, &fsseasonweekid, &value); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("PlayerID : [%d], Value : [%d]\n", playerid, value)

		div := dbfantasyplayervalue.FantasyPlayerValue{
			PlayerID: playerid,
			WeekID:   fsseasonweekid,
			Value:    value,
			FantasyGameID: 2,
		}

		fmt.Printf("Record : %#v\n", div)
		err := dbfantasyplayervalue.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
