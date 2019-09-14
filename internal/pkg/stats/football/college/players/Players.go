package players

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbplayer"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbposition"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbseason"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbsportlevel"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteam"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteamconferenceseason"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
)

var PlayersURL = "https://api.collegefootballdata.com/roster"

type Player struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Weight      int64  `json:"weight"`
	Height      int64  `json:"height"`
	Jersey      int64  `json:"jersey"`
	Year        int64  `json:"year"`
	Position    string `json:"position"`
	HomeCity    string `json:"home_city"`
	HomeState   string `json:"home_state"`
	HomeCountry string `json:"home_country"`
}

func RetrieveAllPlayers() {
	// get sport level
	sportlevel, err := dbsportlevel.ReadByID(dbsportlevel.COLLEGE_FOOTBALL)
	if err != nil {
		log.Fatalf("Error getting sport level record : %v", err)
	}

	// get current season record
	season, err := dbseason.ReadCurrentSeason(sportlevel.CurrentYear, sportlevel.SportLevelID)
	if err != nil {
		log.Fatalf("Error getting current season : %v", err)
	}

	// get all teams in the current season
	teams, err := dbteamconferenceseason.ReadBySeasonID(season.SeasonID)
	if err != nil {
		log.Fatalf("Error getting current all teams : %v", err)
	}

	for _, teamSeason := range teams {

		team, err := dbteam.ReadByID(teamSeason.TeamID)
		if err != nil {
			log.Fatalf("Error getting team : %v", err)
		}

		count := RetrievePlayers(team.Abbreviation.String)
		if count < 1 {
			count = RetrievePlayers(team.Name)

			if count < 1 {
				log.Fatalf("No players were found for %s", team.Name)
			}
		}
	}

}

func RetrievePlayers(teamName string) int {
	count := 0
	url := PlayersURL + "?team=" + url2.QueryEscape(teamName)
	fmt.Printf("Retrieving player data for %s from %s", teamName, url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var players []Player
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bodyBytes)

		json.Unmarshal(bodyBytes, &players)

		fmt.Printf("Players : %+v \n", players)

		// get sport level
		sportlevel, err := dbsportlevel.ReadByID(dbsportlevel.COLLEGE_FOOTBALL)
		if err != nil {
			log.Fatalf("Error getting sport level record : %v", err)
		}

		fmt.Printf("SportLevel : %#v", sportlevel)

		// get team record
		team, err := dbteam.ReadByAbbreviationAndLevel(teamName, dbsportlevel.LEVEL_AMATEUR)
		if err != nil {
			if err == sql.ErrNoRows {
				// try for name
				team, err = dbteam.ReadByNameAndLevel(teamName, dbsportlevel.LEVEL_AMATEUR)
				if err != nil {
					log.Fatalf("Error getting player's team %s with level %s : %v", teamName, dbsportlevel.LEVEL_AMATEUR, err)
				}
			} else {
				log.Fatalf("Error getting player's team %s with level %s : %v", teamName, dbsportlevel.LEVEL_AMATEUR, err)
			}
		}

		for _, playerObj := range players {
			fmt.Println("=========")
			fmt.Printf("Player : %#v \n", playerObj)
			count++

			// get positionObj
			position, err := dbposition.ReadByNameAndSportID(playerObj.Position, sportlevel.SportID)
			if err != nil {
				fmt.Printf("No position object for %s", position)
				continue
			}

			// get player record or create new
			player, err := dbplayer.ReadByStatsKeyAndSportLevel(playerObj.Id, sportlevel.SportLevelID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Fatalf("ERROR: %v", err)
				}

				// create a new
				player = &dbplayer.Player{
					SportLevelID: sportlevel.SportLevelID,
					TeamID:       database.ToNullInt(team.TeamID, true),
					PositionID:   database.ToNullInt(position.PositionID, true),
					FirstName:    database.ToNullString(playerObj.FirstName, true),
					LastName:     database.ToNullString(playerObj.LastName, true),
					Status:       dbplayer.STATUS_ACTIVE,
					StatsKey:     database.ToNullString(playerObj.Id, true),
					Weight:       database.ToNullInt(playerObj.Weight, true),
					Height:       database.ToNullInt(playerObj.Height, true),
					Jersey:       database.ToNullInt(playerObj.Jersey, true),
					Year:         database.ToNullInt(playerObj.Year, true),
					HomeCity:     database.ToNullString(playerObj.HomeCity, true),
					HomeState:    database.ToNullString(playerObj.HomeState, true),
					HomeCountry:  database.ToNullString(playerObj.HomeCountry, true),
				}

				err := dbplayer.Insert(player)
				if err != nil {
					log.Fatalf("ERROR: Can't create ConferenceSeason %v", err)
				}
			}

		}
	}

	return count

}
