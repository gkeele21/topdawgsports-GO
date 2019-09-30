package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbfantasyleague"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "sumo:password@tcp(127.0.0.1:3307)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fsl.FSLeagueID, fsl.LeagueName, fsl.LeaguePassword, fsl.Description, fsl.IsPublic, fss.SeasonID, fss.FSGameID " +
		"FROM fsleague fsl INNER JOIN fsseason fss on fss.FSSeasonID = fsl.FSSeasonID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsleagueid, seasonid, fsgameid int64
		var name, password, desc, ispublic database.NullString
		if err := rows.Scan(&fsleagueid, &name, &password, &desc, &ispublic, &seasonid, &fsgameid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("LeagueID : [%d], Name : [%s]\n", fsleagueid, name)

		league := dbfantasyleague.FantasyLeague{
			FantasyLeagueID: fsleagueid,
			SeasonID:        seasonid,
			FantasyGameID:   fsgameid,
			Name:            name.String,
			Description:     desc,
			Status:          "final",
			LeaguePassword:  password,
			CreatedDate:     time.Now(),
			CreatedByUserID: 1,
		}

		if ispublic.Valid && ispublic.String == "0" {
			league.Visibility = "private"
		} else {
			league.Visibility = "public"
		}

		fmt.Printf("League : %#v\n", league)
		err := dbfantasyleague.Insert(&league)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
