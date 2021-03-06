package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
	"topdawgsportsAPI/pkg/database/dbfantasywaiverwirerequest"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSTeamID, FSSeasonWeekID, Rank, RequestDate, DropPlayerID, DropType, PUPlayerID, PUType, Processed, Granted  FROM FSFootballTransactionRequest fs " +
		"INNER JOIN Player p on p.playerId = fs.DropPlayerId " +
		"INNER JOIN Player p2 on p2.playerId = fs.PUPlayerId ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, fsseasonweekid, granted, processed, rank int64
		var requestdate time.Time
		var dropplayerid, puplayerid database.NullInt64
		var droptype, putype database.NullString
		if err := rows.Scan(&fsteamid, &fsseasonweekid, &rank, &requestdate, &dropplayerid, &droptype, &puplayerid, &putype, &processed, &granted); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], FSSeasonWeekID : [%d], PickupType : [%d]\n", fsteamid, fsseasonweekid, putype)
		if putype.String == "PU" {
			fmt.Println("Changing PickupType to 'pickup'")
			putype.String = "pickup"
		}

		div := dbfantasywaiverwirerequest.FantasyWaiverWireRequest{
			FantasyTeamID:  fsteamid,
			WeekID:         fsseasonweekid,
			RequestDate:    requestdate,
			DropPlayerID:   dropplayerid,
			DropType:       database.ToNullString(strings.ToLower(droptype.String), false),
			PickupPlayerID: puplayerid,
			PickupType:     database.ToNullString(strings.ToLower(putype.String), false),
			Processed:      processed,
			Granted:        granted,
		}

		fmt.Printf("Record : %#v\n", div)
		err := dbfantasywaiverwirerequest.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
