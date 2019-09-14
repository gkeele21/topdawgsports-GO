package games

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbmatchup"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbseason"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbsportlevel"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteam"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbweek"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var GamesURL = "https://api.collegefootballdata.com/games?seasonType=regular"

type Game struct {
	Id             int64               `json:"id"`
	Season         int64               `json:"season"`
	Week           int64               `json:"week"`
	SeasonType     string              `json:"season_type"`
	StartDate      time.Time           `json:"start_date"`
	NeutralSite    bool                `json:"neutral_site"`
	ConferenceGame bool                `json:"conference_game"`
	Attendance     database.NullString `json:"attendance"`
	VenueId        int64               `json:"venue_id"`
	Venue          string              `json:"venue"`
	HomeTeam       string              `json:"home_team"`
	HomeConference string              `json:"home_conference"`
	HomePoints     int64               `json:"home_points"`
	AwayTeam       string              `json:"away_team"`
	AwayConference string              `json:"away_conference"`
	AwayPoints     int64               `json:"away_points"`
}

func RetrieveGamesForWeek(year, week string) []int64 {
	url := GamesURL + "&year=" + year + "&week=" + week
	fmt.Printf("Retrieving game data for year %s and week %s from %s", year, week, url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var gameIds []int64
	var games []Game
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(bodyBytes, &games)

		for _, gameObj := range games {
			fmt.Println("=========")
			fmt.Printf("Game : %#v \n", gameObj)

			gameIds = append(gameIds, gameObj.Id)
		}
	}

	return gameIds
}

func RetrieveGames(year string) int {
	count := 0
	url := GamesURL + "&year=" + year
	fmt.Printf("Retrieving game data for %s from %s", year, url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var games []Game
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

		// get current season record
		season, err := dbseason.ReadCurrentSeason(sportlevel.CurrentYear, sportlevel.SportLevelID)
		if err != nil {
			log.Fatalf("Error getting current season : %v", err)
		}

		for _, gameObj := range games {
			fmt.Println("=========")
			fmt.Printf("Game : %#v \n", gameObj)
			count++

			// get week
			week, err := dbweek.ReadByWeekNameAndSeasonID(strconv.FormatInt(gameObj.Week, 10), season.SeasonID)
			if err != nil {
				if err == sql.ErrNoRows {
					// create new
					week = &dbweek.Week{
						SeasonID:  season.SeasonID,
						WeekName:  strconv.FormatInt(gameObj.Week, 10),
						StartDate: database.NullTime{},
						EndDate:   database.NullTime{},
						Status:    dbweek.STATUS_PENDING,
						WeekType:  dbweek.TYPE_MIDDLE,
					}

					err = dbweek.Insert(week)
					if err != nil {
						log.Fatalf("ERROR: Can't create week %v", err)
					}
				}
				fmt.Printf("No week object for %s", gameObj.Week)
				continue
			}

			// get home team
			homeTeam, err := dbteam.ReadByExternalNameAndLevel(gameObj.HomeTeam, dbsportlevel.LEVEL_AMATEUR)
			if err != nil {
				//if err == sql.ErrNoRows {
				//	// try the name of the team
				//	homeTeam, err = dbteam.ReadByNameAndLevel(gameObj.HomeTeam, dbsportlevel.LEVEL_AMATEUR)
				//	if err != nil {
				//		log.Fatalf("ERROR: Can't find home team with Name : %s - %s", gameObj.HomeTeam, err)
				//	}
				//} else {
				log.Fatalf("ERROR: Can't find home team with ExternalName : %s - %s", gameObj.HomeTeam, err)
				//}
			}

			// get away team
			awayTeam, err := dbteam.ReadByExternalNameAndLevel(gameObj.AwayTeam, dbsportlevel.LEVEL_AMATEUR)
			if err != nil {
				//if err == sql.ErrNoRows {
				//	// try the name of the team
				//	awayTeam, err = dbteam.ReadByNameAndLevel(gameObj.AwayTeam, dbsportlevel.LEVEL_AMATEUR)
				//	if err != nil {
				//		log.Fatalf("ERROR: Can't find away team with Name : %s", gameObj.AwayTeam)
				//	}
				//} else {
				log.Fatalf("ERROR: Can't find away team with ExternalName : %s", gameObj.AwayTeam)
				//}
			}

			// get matchup record or create new
			matchup, err := dbmatchup.ReadByWeekIDAndVisitorIDAndHomeId(week.WeekID, awayTeam.TeamID, homeTeam.TeamID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Fatalf("ERROR: %v", err)
				}

				// create a new
				matchup = &dbmatchup.Matchup{
					WeekID:        database.ToNullInt(week.WeekID, true),
					MatchupDate:   database.ToNullTime(gameObj.StartDate, true),
					VisitorTeamID: database.ToNullInt(awayTeam.TeamID, true),
					HomeTeamID:    database.ToNullInt(homeTeam.TeamID, true),
					VenueID:       database.NullInt64{},
					VisitorScore:  database.ToNullInt(gameObj.AwayPoints, true),
					HomeScore:     database.ToNullInt(gameObj.HomePoints, true),
					WinningTeamID: database.NullInt64{},
					NumOvertimes:  database.NullInt64{},
					Status:        dbmatchup.STATUS_PENDING,
					ExternalId:    database.ToNullInt(gameObj.Id, true),
				}

				matchup.SetWinner()

				err := dbmatchup.Insert(matchup)
				if err != nil {
					log.Fatalf("ERROR: Can't create Matchup %v", err)
				}
			} else {
				// update score
				matchup.HomeScore = database.ToNullInt(gameObj.HomePoints, true)
				matchup.VisitorScore = database.ToNullInt(gameObj.AwayPoints, true)

				matchup.SetWinner()
				if !matchup.ExternalId.Valid {
					matchup.ExternalId = database.ToNullInt(gameObj.Id, true)
				}

				dbmatchup.Update(matchup)
			}

		}
	}

	return count
}
