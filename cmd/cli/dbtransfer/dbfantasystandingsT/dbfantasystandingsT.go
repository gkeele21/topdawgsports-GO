package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasyheadtoheadstandings"
	"topdawgsportsAPI/pkg/database/dbfantasysalarystandings"
	"topdawgsportsAPI/pkg/database/dbloveemstandings"
	"topdawgsportsAPI/pkg/database/dbpickemstandings"
)

var db *sql.DB

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fs.FSTeamID, fs.FSSeasonWeekID, FantasyPts, TotalFantasyPts, GamePoints, TotalGamePoints, Wins, Losses, Ties," +
		" FantasyPtsAgainst, HiScore, TotalHiScores, CurrentStreak, LastFive, Rank, TotalFantasyPtsAgainst, GamesCorrect, TotalGamesCorrect, GamesWrong, TotalGamesWrong, g.FSGameID " +
		" FROM FSFootballStandings fs" +
		" INNER JOIN FSTeam t on t.FSTeamID = fs.FSTeamID " +
		"INNER JOIN FSLeague l on l.FSLeagueID = t.FSLeagueID " +
		"INNER JOIN FSSeason se on se.FSSeasonID = l.FSSeasonID " +
		"INNER JOIN FSGame g on g.FSGameID = se.FSGameID " +
		"INNER JOIN FSSeasonWeek fsw on fsw.FSSeasonWeekID = fs.FSSeasonWeekID ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fsteamid, fsseasonweekid, fsgameid int64
		var fantasypts, totalfantasypts, fantasyptsagainst, totalfantasyptsagainst, gamepoints, totalgamepoints database.NullFloat64
		var wins, losses, ties, hiscore, totalhiscores, currentstreak, rank, gamescorrect, totalgamescorrect, gameswrong, totalgameswrong database.NullInt64
		var lastfive database.NullString
		if err := rows.Scan(&fsteamid, &fsseasonweekid, &fantasypts, &totalfantasypts, &gamepoints, &totalgamepoints, &wins, &losses, &ties, &fantasyptsagainst, &hiscore, &totalhiscores, &currentstreak, &lastfive, &rank, &totalfantasyptsagainst, &gamescorrect, &totalgamescorrect, &gameswrong, &totalgameswrong, &fsgameid); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("FSTeamID : [%d], FSSeasonWeekID : [%d]\n", fsteamid, fsseasonweekid)

		if fsgameid == 1 {
			div := dbfantasyheadtoheadstandings.FantasyHeadToHeadStandings{
				FantasyTeamID:          fsteamid,
				WeekID:                 fsseasonweekid,
				WeekFantasyPts:         fantasypts,
				TotalFantasyPts:        totalfantasypts,
				WeekFantasyPtsAgainst:  fantasyptsagainst,
				TotalFantasyPtsAgainst: totalfantasyptsagainst,
				Wins:                   wins.Int64,
				Losses:                 losses.Int64,
				Ties:                   ties.Int64,
				WeekHiScore:            hiscore,
				TotalHiScore:           totalhiscores,
				LeagueRanking:          rank,
				CurrentStreak:          currentstreak,
				LastFive:               lastfive,
			}

			fmt.Printf("Record : %#v\n", div)
			err := dbfantasyheadtoheadstandings.Insert(&div)
			if err != nil {
				log.Fatal(err)
			}

		} else if fsgameid == 4 || fsgameid == 7 {
			div := dbpickemstandings.PickEmStandings{
				FantasyTeamID:     fsteamid,
				WeekID:            fsseasonweekid,
				WeekGamePoints:    int64(gamepoints.Float64),
				TotalGamePoints:   int64(totalgamepoints.Float64),
				LeagueRanking:     rank,
				GamesCorrect:      gamescorrect,
				TotalGamesCorrect: totalgamescorrect,
				GamesWrong:        gameswrong,
				TotalGamesWrong:   totalgameswrong,
			}

			fmt.Printf("Record : %#v\n", div)
			err := dbpickemstandings.Insert(&div)
			if err != nil {
				log.Fatal(err)
			}

		} else if fsgameid == 2 {
			div := dbfantasysalarystandings.FantasySalaryStandings{
				FantasyTeamID:   fsteamid,
				WeekID:          fsseasonweekid,
				WeekFantasyPts:  fantasypts.Float64,
				TotalFantasyPts: totalfantasypts.Float64,
				WeekGamePts:     int64(gamepoints.Float64),
				TotalGamePts:    int64(totalgamepoints.Float64),
				LeagueRanking:   rank,
			}

			fmt.Printf("Record : %#v\n", div)
			err := dbfantasysalarystandings.Insert(&div)
			if err != nil {
				log.Fatal(err)
			}

		} else if fsgameid == 5 || fsgameid == 6 {
			div := dbloveemstandings.LoveEmStandings{
				FantasyTeamID: fsteamid,
				WeekID:        fsseasonweekid,
				WeekGamePts:   int64(gamepoints.Float64),
				TotalGamePts:  int64(totalgamepoints.Float64),
				LeagueRanking: rank,
			}

			fmt.Printf("Record : %#v\n", div)
			err := dbloveemstandings.Insert(&div)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Print("No record to insert")
			os.Exit(0)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
