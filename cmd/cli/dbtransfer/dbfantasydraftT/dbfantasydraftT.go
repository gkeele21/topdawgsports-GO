package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasydraft"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fd.FSLeagueID, fd.Round, fd.Place, fd.FSTeamID, fd.PlayerID FROM FSFootballDraft fd INNER JOIN Player p on p.PlayerID = fd.PlayerID INNER JOIN FSTeam ft on ft.FSTeamID = fd.FSTeamID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsleagueid, round, place int64
		var fsteamid, playerid database.NullInt64
		if err := rows.Scan(&fsleagueid, &round, &place, &fsteamid, &playerid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Round: [%d], Place : [%d]\n", round, place)

		draft := dbfantasydraft.FantasyDraft{
			FantasyLeagueID: fsleagueid,
			Round:           round,
			Place:           place,
			FantasyTeamID:   fsteamid,
			PlayerID:        playerid,
		}

		fmt.Printf("Draft : %#v\n", draft)
		err := dbfantasydraft.Insert(&draft)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
