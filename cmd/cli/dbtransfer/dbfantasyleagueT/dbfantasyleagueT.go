package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasyleague"
	"time"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fsl.FSLeagueID, fsl.LeagueName, fsl.LeaguePassword, fsl.Description, fsl.IsPublic, fss.SeasonID, fss.FSGameID " +
		"FROM FSLeague fsl INNER JOIN FSSeason fss on fss.FSSeasonID = fsl.FSSeasonID")
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
