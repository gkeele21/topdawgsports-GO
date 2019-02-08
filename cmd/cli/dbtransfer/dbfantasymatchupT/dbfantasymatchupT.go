package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasymatchup"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fm.FSLeagueID, fm.FSSeasonWeekID, fm.Team1ID, fm.Team2ID, fm.Team1Pts, fm.Team2Pts, fm.Winner FROM FSFootballMatchup fm")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsleagueid, fsseasonweekid int64
		var team1id, team2id, winner database.NullInt64
		var team1pts, team2pts database.NullFloat64
		if err := rows.Scan(&fsleagueid, &fsseasonweekid, &team1id, &team2id, &team1pts, &team2pts, &winner); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSLeagueID : [%d], FSSeasonWeekID: [%d]\n", fsleagueid, fsseasonweekid)

		div := dbfantasymatchup.FantasyMatchup{
			FantasyLeagueID: fsleagueid,
			WeekID:          fsseasonweekid,
			VisitorTeamID:   team1id,
			VisitorScore:    team1pts,
			HomeTeamID:      team2id,
			HomeScore:       team2pts,
		}
		if winner.Int64 > 0 {
			div.WinningTeamID = winner
		}

		fmt.Printf("Matchup : %#v\n", div)
		err := dbfantasymatchup.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
