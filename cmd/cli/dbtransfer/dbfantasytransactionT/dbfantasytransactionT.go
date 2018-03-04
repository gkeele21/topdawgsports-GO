package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasytransaction"
	"time"
	"strings"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSTeamID, FSSeasonWeekID, FSLeagueID, TransactionDate, DropPlayerID, DropType, PUPlayerID, PUType  FROM FSFootballTransaction fs " +
		"INNER JOIN Player p on p.playerId = fs.DropPlayerId " +
		"INNER JOIN Player p2 on p2.playerId = fs.PUPlayerId ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, fsseasonweekid, fsleagueid int64
		var transactiondate time.Time
		var dropplayerid, puplayerid database.NullInt64
		var droptype, putype database.NullString
		if err := rows.Scan(&fsteamid, &fsseasonweekid, &fsleagueid, &transactiondate, &dropplayerid, &droptype, &puplayerid, &putype); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], FSSeasonWeekID : [%d], PickupType : [%d]\n", fsteamid, fsseasonweekid, putype)
		if putype.String == "PU" {
			fmt.Println("Changing PickupType to 'pickup'")
			putype.String = "pickup"
		}

		div := dbfantasytransaction.FantasyTransaction{
			FantasyLeagueID: fsleagueid,
			FantasyTeamID:   fsteamid,
			WeekID:          database.ToNullInt(fsseasonweekid, false),
			TransactionDate: transactiondate,
			DropPlayerID:    dropplayerid,
			DropType:        database.ToNullString(strings.ToLower(droptype.String), false),
			PickupPlayerID:  puplayerid,
			PickupType:      database.ToNullString(strings.ToLower(putype.String), false),
		}

		fmt.Printf("FantasyTransaction : %#v\n", div)
		err := dbfantasytransaction.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
