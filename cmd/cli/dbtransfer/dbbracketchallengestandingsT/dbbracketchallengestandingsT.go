package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbbracketchallengestandings"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSTeamID, FSSeasonWeekID, Rank, RoundPoints, TotalPoints, MaxPossible FROM BracketChallengeStandings")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, fsseasonweekid, roundpoints, totalpoints, maxpossible int64
		var rank database.NullInt64
		if err := rows.Scan(&fsteamid, &fsseasonweekid, &rank, &roundpoints, &totalpoints, &maxpossible); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], FSSeasonWeekID : [%d]\n", fsteamid, fsseasonweekid)

		div := dbbracketchallengestandings.BracketChallengeStandings{
			FantasyTeamID: fsteamid,
			WeekID:        fsseasonweekid,
			Rank:          rank,
			RoundPoints:   roundpoints,
			TotalPoints:   totalpoints,
			MaxPossible:   database.ToNullInt(maxpossible, false),
		}

		fmt.Printf("Record : %v\n", div)
		err := dbbracketchallengestandings.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
