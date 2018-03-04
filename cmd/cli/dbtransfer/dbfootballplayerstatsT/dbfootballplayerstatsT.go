package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"topdawgsportsAPI/pkg/database/dbfootballplayerstats"
	"topdawgsportsAPI/pkg/database"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT StatsPlayerID, SeasonID, SeasonWeekID, Started, Played, Inactive, PassComp, PassAtt, PassYards, PassInt, PassTD, PassTwoPt, Sacked, SackedYardsLost, RushAtt, RushYards, RushTD, RushTwoPt, RecTargets, RecCatches, RecYards, RecTD, RecTwoPt, XPM, XPA, XPBlocked, FGM, FGA, FG29Minus, FG30to39, FG40to49, FG50Plus, FumblesLost, XtraTD, TDDistances FROM FootballStats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var seasonid int64
		var weekid, started, played, inactive, passatt, passcomp database.NullInt64
		var passyards, passint, passtds, passtwopt, sacked, sackedyardslost, rushatt, rushyards database.NullInt64
		var rushtd, rushtwopt, rectargets, reccatches, recyards, rectd, rectwopt, xpm database.NullInt64
		var xpa, xpblocked, fgm, fga, fg29, fg30, fg40, fg50, fumbleslost, xtratd database.NullInt64
		var statsplayerid, tddistances database.NullString
		if err := rows.Scan(&statsplayerid, &seasonid, &weekid, &started, &played, &inactive, &passcomp, &passatt, &passyards, &passint, &passtds, &passtwopt, &sacked, &sackedyardslost, &rushatt, &rushyards, &rushtd, &rushtwopt, &reccatches, &rectargets, &recyards, &rectd, &rectwopt, &xpm, &xpa, &xpblocked, &fgm, &fga, &fg29, &fg30, &fg40, &fg50, &fumbleslost, &xtratd, &tddistances); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("SeasonID : [%d], StatsPlayerID : [%d], WeekID : [%d]\n", seasonid, statsplayerid, weekid)

		// grab the PlayerID based on the StatsID using NFLGameStatsID or StatsPlayerID
		var playerid int64
		row := db.QueryRow("SELECT PlayerID FROM Player WHERE StatsPlayerID = ?", statsplayerid)

		fmt.Printf("StatsPlayerID : [%s]\n", statsplayerid)
		switch err := row.Scan(&playerid); err {
		case nil:
			fmt.Println(playerid)
		default:
			// try to get look up the player by the nflgamestatsid
			row2 := db.QueryRow("SELECT PlayerID FROM Plasyer WHERE NFLGameStatsID = ?", statsplayerid)

			switch err := row2.Scan(&playerid); err {
			case nil:
				fmt.Println(playerid)
			default:

			}

		}

		if playerid > 0 {
			stat := dbfootballplayerstats.FootballPlayerStats{
				SeasonID:          seasonid,
				PlayerID:          playerid,
				Started:           started,
				Played:            played,
				Inactive:          inactive,
				PassComp:          passcomp,
				PassAttempts:      passatt,
				PassYards:         passyards,
				PassInterceptions: passint,
				PassTDs:           passtds,
				PassTwoPts:        passtwopt,
				Sacked:            sacked,
				SackedYardsLost:   sackedyardslost,
				RushAttempts:      rushatt,
				RushYards:         rushyards,
				RushTDs:           rushtd,
				RushTwoPts:        rushtwopt,
				RecTargets:        rectargets,
				RecCatches:        reccatches,
				RecYards:          recyards,
				RecTDs:            rectd,
				RecTwoPts:         rectwopt,
				XPMade:            xpm,
				XPAttempts:        xpa,
				XPBlocked:         xpblocked,
				FGMade:            fgm,
				FGAttempts:        fga,
				FG29Minus:         fg29,
				FG3039:            fg30,
				FG4049:            fg40,
				FG50Plus:          fg50,
				FumblesLost:       fumbleslost,
				ExtraTDs:          xtratd,
				TDDistances:       tddistances,
			}
			if weekid.Int64 > 0 {
				stat.WeekID = weekid
			}

			fmt.Printf("Stats : %#v\n", stat)
			err := dbfootballplayerstats.Insert(&stat)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
