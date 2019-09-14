package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbpickempick"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT c.FSSeasonWeekID, c.FSTeamID, c.GameID, c.TeamPickedID, c.ConfidencePts FROM CollegePickem c INNER JOIN Game g on g.GameID = c.GameID INNER JOIN FSTeam t on t.FSTeamID = c.FSTeamID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsseasonweekid, fsteamid, gameid, teampickedid int64
		var confidencepts database.NullInt64
		if err := rows.Scan(&fsseasonweekid, &fsteamid, &gameid, &teampickedid, &confidencepts); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSSeasonWeekID : [%d], FSTeamID : [%d], GameID : [%d]\n", fsseasonweekid, fsteamid, gameid)

		div := dbpickempick.PickEmPick{
			FantasyTeamID: fsteamid,
			WeekID:        fsseasonweekid,
			MatchupID:     gameid,
			TeamPickedID:  teampickedid,
		}

		if confidencepts.Int64 > 0 {
			div.ConfidencePoints = confidencepts
		}

		fmt.Printf("Pick : %#v\n", div)
		err := dbpickempick.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
