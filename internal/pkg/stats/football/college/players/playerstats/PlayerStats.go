package playerstats

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbfootballplayerstats"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbmatchup"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbplayer"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbsportlevel"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteam"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbweek"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var GameStatsURL = "https://api.collegefootballdata.com/games/players?seasonType=regular"

type GameStats struct {
	Id    int         `json:"id"`
	Teams []TeamStats `json:"teams"`
}

type TeamStats struct {
	School     string       `json:"school"`
	Conference string       `json:"conference"`
	HomeAway   string       `json:"homeAway"`
	Points     int          `json:"points"`
	Categories []GroupStats `json:"categories"`
}

type GroupStats struct {
	Name  string      `json:"name"`
	Types []TypeStats `json:"types"`
}

type TypeStats struct {
	Name     string        `json:"name"`
	Athletes []PlayerStats `json:"athletes"`
}

type PlayerStats struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Stat string `json:"stat"`
}

func RetrieveGamePlayerStats(year string, gameId int64) {
	// get Game from db
	matchup, err := dbmatchup.ReadByExternalID(gameId)
	if err != nil {
		log.Fatalf("ERROR: Invalid GameId passed : %s", gameId)
	}

	weekId := matchup.WeekID
	week, err := dbweek.ReadByID(weekId.Int64)
	if err != nil {
		log.Fatalf("ERROR: The matchup has an invalid weekId")
	}

	seasonId := week.SeasonID

	url := GameStatsURL + "?year=" + year + "&gameId=" + strconv.FormatInt(gameId, 10)
	fmt.Printf("Retrieving player stats data for gameId %s from %s", gameId, url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var games []GameStats
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bodyBytes)

		json.Unmarshal(bodyBytes, &games)

		fmt.Printf("Games : %+v \n", games)

		// get sport level
		sportlevel, err := dbsportlevel.ReadByID(dbsportlevel.COLLEGE_FOOTBALL)
		if err != nil {
			log.Fatalf("Error getting sport level record : %v", err)
		}

		for _, gameObj := range games {
			fmt.Println("=========")
			fmt.Printf("Game : %#v \n", gameObj)

			for _, teamObj := range gameObj.Teams {
				fmt.Println("------")
				fmt.Printf("Team : %#v \n", teamObj)

				// get team from db
				team, err := dbteam.ReadByExternalNameAndLevel(teamObj.School, dbsportlevel.LEVEL_AMATEUR)
				if err != nil {
					log.Fatalf("ERROR: Can't find away team with ExternalName : %s", teamObj.School)
				}
				fmt.Printf("TeamId : %v \n", team.TeamID)

				for _, category := range teamObj.Categories {
					fmt.Println("------")
					fmt.Printf("Category : %#v \n", category)

					processCategory(seasonId, weekId.Int64, sportlevel.SportLevelID, category, *team)

				}
			}
		}
	}
}

func processCategory(seasonId, weekId, sportLevelId int64, category GroupStats, team dbteam.Team) {
	fmt.Println("------")
	fmt.Printf("Category : %#v \n", category)

	for _, categoryType := range category.Types {
		fmt.Printf("CategoryType : %#v \n", categoryType)

		for _, athlete := range categoryType.Athletes {
			categoryName := category.Name

			if categoryName == "defensive" || categoryName == "punting" || categoryName == "interceptions" {
				continue
			}

			// get Player Obj
			player, err := dbplayer.ReadByStatsKeyAndSportLevel(athlete.Id, sportLevelId)
			if err != nil {
				fmt.Printf("ERROR: could not find a player with the statId of %s - Name = %s - %s", athlete.Id, athlete.Name, err)
				continue
			}

			statsRec := getPlayerStatsRecord(seasonId, weekId, player.PlayerID)

			switch categoryName {
			case "defensive":
				// ignore for now
			case "fumbles":
				handleFumbles(categoryType.Name, statsRec, athlete.Stat)
			case "punting":
				// ignore for now
			case "kicking":
				handleKicking(categoryType.Name, statsRec, athlete.Stat)
			case "puntReturns":
				handlePuntReturns(categoryType.Name, statsRec, athlete.Stat)
			case "kickReturns":
				handleKickReturns(categoryType.Name, statsRec, athlete.Stat)
			case "interceptions":
				// ignore for now
			case "receiving":
				handleReceiving(categoryType.Name, statsRec, athlete.Stat)
			case "rushing":
				handleRushing(categoryType.Name, statsRec, athlete.Stat)
			case "passing":
				handlePassing(categoryType.Name, statsRec, athlete.Stat)
			}
		}
	}
}

func getPlayerStatsRecord(seasonId, weekId, playerId int64) dbfootballplayerstats.FootballPlayerStats {
	// get or create Player stats record
	statsRec, err := dbfootballplayerstats.ReadBySeasonWeekPlayer(seasonId, weekId, playerId)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("ERROR: retrieving stats record : %s", err)
		}
		statsRec = &dbfootballplayerstats.FootballPlayerStats{
			SeasonID: seasonId,
			WeekID:   database.ToNullInt(weekId, true),
			PlayerID: playerId,
		}
	}

	return *statsRec
}

func getStatIntValue(stat string) database.NullInt64 {
	statValue, err := strconv.ParseInt(stat, 10, 64)
	if err != nil {
		statValue = 0
	}

	return database.ToNullInt(statValue, true)
}

func handleFumbles(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "REC":
		statsRec.DefFumbleRecoveries = getStatIntValue(statValue)
	case "LOST":
		statsRec.FumblesLost = getStatIntValue(statValue)
	case "FUM":
		statsRec.Fumbles = getStatIntValue(statValue)
	}
	dbfootballplayerstats.Save(&statsRec)
}

func handleKicking(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "XP":
		fmt.Printf("C/ATT orig: %s\n", statValue)
		// parse out the c/att
		stringArr := strings.Split(statValue, "/")
		fmt.Print("C/ATT: ")
		fmt.Println(stringArr)
		intVal, err := strconv.ParseInt(stringArr[0], 10, 64)
		if err != nil {
			intVal = 0
		}
		fmt.Printf("Setting XP made value to %#v \n", database.ToNullInt(intVal, true))
		statsRec.XPMade = database.ToNullInt(intVal, true)
		fmt.Printf("XPMade value 230 : %#v \n", statsRec.XPMade)

		intVal, err = strconv.ParseInt(stringArr[1], 10, 64)
		if err != nil {
			intVal = 0
		}
		statsRec.XPAttempts = database.ToNullInt(intVal, true)
	case "FG":
		// we only get the # of field goals made - no distances
		fmt.Printf("C/ATT orig: %s\n", statValue)
		// parse out the c/att
		stringArr := strings.Split(statValue, "/")
		fmt.Print("C/ATT: ")
		fmt.Println(stringArr)
		intVal, err := strconv.ParseInt(stringArr[0], 10, 64)
		if err != nil {
			intVal = 0
		}
		statsRec.FGMade = database.ToNullInt(intVal, true)

		intVal, err = strconv.ParseInt(stringArr[1], 10, 64)
		if err != nil {
			intVal = 0
		}
		statsRec.FGAttempts = database.ToNullInt(intVal, true)
	}
	dbfootballplayerstats.Save(&statsRec)
}

func handlePuntReturns(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "TD":
		statsRec.ExtraTDs = database.ToNullInt(statsRec.ExtraTDs.Int64+getStatIntValue(statValue).Int64, true)
	}
	dbfootballplayerstats.Save(&statsRec)
}

func handleKickReturns(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "TD":
		statsRec.ExtraTDs = database.ToNullInt(statsRec.ExtraTDs.Int64+getStatIntValue(statValue).Int64, true)
	}
	dbfootballplayerstats.Save(&statsRec)
}

func handlePassing(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "INT":
		statsRec.PassInterceptions = getStatIntValue(statValue)
	case "TD":
		statsRec.PassTDs = getStatIntValue(statValue)
	case "YDS":
		statsRec.PassYards = getStatIntValue(statValue)
	case "C/ATT":
		fmt.Printf("C/ATT orig: %s\n", statValue)
		// parse out the c/att
		stringArr := strings.Split(statValue, "/")
		fmt.Print("C/ATT: ")
		fmt.Println(stringArr)
		intVal, err := strconv.ParseInt(stringArr[0], 10, 64)
		if err != nil {
			intVal = 0
		}
		statsRec.PassComp = database.ToNullInt(intVal, true)

		intVal, err = strconv.ParseInt(stringArr[1], 10, 64)
		if err != nil {
			intVal = 0
		}
		statsRec.PassAttempts = database.ToNullInt(intVal, true)

	}
	dbfootballplayerstats.Save(&statsRec)
}

func handleReceiving(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "TD":
		statsRec.RecTDs = getStatIntValue(statValue)
	case "YDS":
		statsRec.RecYards = getStatIntValue(statValue)
	case "REC":
		statsRec.RecCatches = getStatIntValue(statValue)
	}
	dbfootballplayerstats.Save(&statsRec)
}

func handleRushing(categoryName string, statsRec dbfootballplayerstats.FootballPlayerStats, statValue string) {
	switch categoryName {
	case "TD":
		statsRec.RushTDs = getStatIntValue(statValue)
	case "YDS":
		statsRec.RushYards = getStatIntValue(statValue)
	case "CAR":
		statsRec.RushAttempts = getStatIntValue(statValue)
	}
	dbfootballplayerstats.Save(&statsRec)
}
