package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbplayer"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

var db *sql.DB

func main() {
	// grab all players from the existing database
	db, err := sql.Open("mysql", "sumo:password@tcp(127.0.0.1:3307)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT PlayerID, TeamID, PositionID, FirstName, LastName, IsActive, NFLGameStatsID, StatsPlayerID FROM player")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var playerid, isactive int64
		var teamid, positionid, statsplayerid database.NullInt64
		var firstname, lastname, statsid database.NullString
		if err := rows.Scan(&playerid, &teamid, &positionid, &firstname, &lastname, &isactive, &statsid, &statsplayerid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("PlayerID : [%d], TeamID : [%d], PositionID : [%d], First : [%s], Last : [%s], Active : [%d], StatsId : [%s]\n", playerid, teamid, positionid, firstname, lastname, isactive, statsid)

		status := "active"
		if isactive == 0 {
			status = "inactive"
		}

		if statsid.String == "null" {
			statsid = database.ToNullString(strconv.FormatInt(statsplayerid.Int64, 10), true)
		}
		player := dbplayer.Player{
			PlayerID:  playerid,
			FirstName: firstname,
			LastName:  lastname,
			Status:    status,
			StatsKey:  statsid,
		}

		sportlevelid := int64(1)
		if teamid.Int64 != 0 {
			player.TeamID = teamid
		}

		if positionid.Int64 > 0 && positionid.Int64 <= 7 {
			player.PositionID = positionid
		}

		if statsid.String != "-1" && statsid.String != "null" && statsid.String != "0" {
			player.StatsKey = statsid
		}

		if positionid.Int64 >= 12 {
			sportlevelid = 4
			player.PositionID = database.ToNullInt(8, true)
		}

		player.SportLevelID = sportlevelid
		fmt.Printf("Player : %v\n", player)
		err := dbplayer.Insert(&player)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
