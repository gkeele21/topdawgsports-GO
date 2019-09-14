package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"topdawgsportsAPI/pkg/database/dbfantasyteam"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fst.FSTeamID, fst.FSLeagueID, fst.FSUserID, fst.TeamName, fst.DateCreated, fst.ScheduleTeamNo FROM FSTeam fst INNER JOIN FSUser u on u.FSUserID = fst.FSUserID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, fsleagueid, fsuserid int64
		var teamname string
		var dateCreated time.Time
		var scheduleteamno database.NullInt64
		if err := rows.Scan(&fsteamid, &fsleagueid, &fsuserid, &teamname, &dateCreated, &scheduleteamno); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], Name : [%s]\n", fsteamid, teamname)

		team := dbfantasyteam.FantasyTeam{
			FantasyTeamID:      fsteamid,
			FantasyLeagueID:    fsleagueid,
			UserID:             fsuserid,
			Name:               teamname,
			Created:            dateCreated,
			ScheduleTeamNumber: scheduleteamno,
			Status:             "active",
		}

		fmt.Printf("Team : %#v\n", team)
		err := dbfantasyteam.Insert(&team)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
