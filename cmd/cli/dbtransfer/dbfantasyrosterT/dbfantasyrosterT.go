package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasyroster"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fr.ID, fr.FSTeamID, fr.FSSeasonWeekID, fr.PlayerID, fr.StarterState, fr.ActiveState FROM FSRoster fr INNER JOIN Player p on p.PlayerID = fr.PlayerID INNER JOIN FSTeam fst on fst.FSTeamID = fr.FSTeamID INNER JOIN FSSeasonWeek fsw ON fsw.FSSeasonWeekID = fr.FSSeasonWeekID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsrosterid, fsteamid, fsseasonweekid int64
		var starterstate, activestate string
		var playerid database.NullInt64
		if err := rows.Scan(&fsrosterid, &fsteamid, &fsseasonweekid, &playerid, &starterstate, &activestate); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSRosterID : %d, FSTeamID : %d, PlayerID : %d, ActiveState : %s, StarterState : %#v, \n", fsrosterid, fsteamid, playerid, activestate, starterstate)

		ros := dbfantasyroster.FantasyRoster{
			FantasyTeamID: fsteamid,
			WeekID:        fsseasonweekid,
			PlayerID:      playerid,
		}

		if activestate == "ir" {
			ros.ScoringState = "ir"
		} else if starterstate == "starter" {
			fmt.Println("setting to starter")
			ros.ScoringState = "scoring"
		} else if starterstate == "bench" {
			ros.ScoringState = "bench"
		} else {
			fmt.Println("No match")
		}

		fmt.Printf("Roster : %#v\n", ros)

		err := dbfantasyroster.Insert(&ros)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
