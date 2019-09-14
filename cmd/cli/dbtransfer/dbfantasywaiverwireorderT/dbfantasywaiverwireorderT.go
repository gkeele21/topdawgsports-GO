package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbfantasywaiverwireorder"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSLeagueID, fto.FSSeasonWeekID, OrderNumber, FSTeamID FROM FSFootballTransactionOrder fto " +
		"INNER JOIN FSSeasonWeek w on w.FSSeasonWeekID = fto.FSSeasonWeekID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsleagueid, fsseasonweekid, ordernumber int64
		var fsteamid database.NullInt64
		if err := rows.Scan(&fsleagueid, &fsseasonweekid, &ordernumber, &fsteamid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSLeagueID : [%d], Week : [%d]\n", fsleagueid, fsseasonweekid)

		div := dbfantasywaiverwireorder.FantasyWaiverWireOrder{
			FantasyLeagueID: fsleagueid,
			WeekID:          fsseasonweekid,
			OrderNumber:     ordernumber,
			FantasyTeamID:   fsteamid,
		}

		fmt.Printf("Record : %#v\n", div)
		err := dbfantasywaiverwireorder.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
