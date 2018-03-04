package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbmatchup"
	"topdawgsportsAPI/pkg/database/dbteam"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	var err error
	db, err = sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT GameID, SeasonWeekID, GameDate, VisitorID, HomeID, VisitorPts, HomePts, WinnerID, NumOTs  FROM Game")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var gameid, seasonweekid int64
		var gamedate database.NullTime
		var tempvisitorid, temphomeid, tempwinnerid, numots, visitorPts, homePts database.NullInt64
		if err := rows.Scan(&gameid, &seasonweekid, &gamedate, &tempvisitorid, &temphomeid, &visitorPts, &homePts, &tempwinnerid, &numots); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("GameID : [%d], SeasonWeekID : [%d]\n", gameid, seasonweekid)

		// grab the TeamID based on the StatsTeamID using LiveStatsTeamID or StatsTeamID
		winnerteam, err := dbteam.ReadByID(tempwinnerid.Int64)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
		var winnerid int64
		if winnerteam != nil && winnerteam.TeamID >= 1 {
			winnerid = winnerteam.TeamID
		} else {
			winnerid = getTeamID(tempwinnerid.Int64)
		}

		visitorteam, err := dbteam.ReadByID(tempvisitorid.Int64)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
		var visitorid int64
		if visitorteam != nil && visitorteam.TeamID >= 1 {
			visitorid = visitorteam.TeamID
		} else {
			visitorid = getTeamID(tempvisitorid.Int64)
		}

		hometeam, err := dbteam.ReadByID(temphomeid.Int64)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}

		var homeid int64
		if hometeam != nil && hometeam.TeamID >= 1 {
			homeid = hometeam.TeamID
		} else {
			homeid = getTeamID(temphomeid.Int64)
		}

		match := dbmatchup.Matchup{
			MatchupID:     gameid,
			WeekID:        database.ToNullInt(seasonweekid, false),
			MatchupDate:   gamedate,
			VisitorTeamID: database.ToNullInt(visitorid, false),
			HomeTeamID:    database.ToNullInt(homeid, false),
			VisitorScore:  visitorPts,
			HomeScore:     homePts,
			WinningTeamID: database.ToNullInt(winnerid, false),
			NumOvertimes:  numots,
			Status:        "final",
		}

		fmt.Printf("Matchup : %#v\n", match)
		err = dbmatchup.Insert(&match)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func getTeamID(dbteamid int64) int64 {
	// grab the TeamID based on the StatsTeamID using LiveStatsTeamID or StatsTeamID
	var teamid int64
	fmt.Printf("LiveStatsTeamID : [%d]\n", dbteamid)
	row := db.QueryRow("SELECT TeamID FROM Team WHERE LiveStatsTeamID = ?", dbteamid)

	switch err := row.Scan(&teamid); err {
	case nil:
		fmt.Println(teamid)
	default:
		// try to get look up the team by the StatsTeamID
		row2 := db.QueryRow("SELECT TeamID FROM Team WHERE StatsTeamID = ?", dbteamid)

		switch err2 := row2.Scan(&teamid); err2 {
		case nil:
			fmt.Println(teamid)
		default:

		}

	}

	return teamid

}
