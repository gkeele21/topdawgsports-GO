package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbbracketchallengepick"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSTeamID, GameID, TeamSeedPickedID FROM BracketChallenge")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, gameid, pickedid int64
		if err := rows.Scan(&fsteamid, &gameid, &pickedid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], GameID: [%d], PickedID: [%d]\n", fsteamid, gameid, pickedid)

		pick := dbbracketchallengepick.BracketChallengePick{
			FantasyTeamID:    fsteamid,
			MatchupID:        gameid,
			TeamSeedPickedID: database.ToNullInt(pickedid, false),
		}

		fmt.Printf("Pick : %v\n", pick)
		err := dbbracketchallengepick.Insert(&pick)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
