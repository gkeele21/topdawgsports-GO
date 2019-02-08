package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbloveempick"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT c.FSSeasonWeekID, c.FSTeamID, c.TeamPickedID FROM CollegeLoveLeave c INNER JOIN FSTeam t on t.FSTeamID = c.FSTeamID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsseasonweekid, fsteamid, teampickedid int64
		if err := rows.Scan(&fsseasonweekid, &fsteamid, &teampickedid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSSeasonWeekID : [%d], FSTeamID : [%d], TeamPickedID : [%d]\n", fsseasonweekid, fsteamid, teampickedid)

		div := dbloveempick.LoveEmPick{
			WeekID:        fsseasonweekid,
			FantasyTeamID: fsteamid,
			TeamPickedID:  teampickedid,
		}

		fmt.Printf("Record : %v\n", div)
		err := dbloveempick.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
