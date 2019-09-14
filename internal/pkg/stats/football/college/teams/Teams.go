package teams

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbconference"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbconferenceseason"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbseason"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbsportlevel"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteam"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbteamconferenceseason"
	"io/ioutil"
	"log"
	"net/http"
)

var TeamsURL = "https://api.collegefootballdata.com/teams/fbs"

type School struct {
	Id           string   `json:"id"`
	Name         string   `json:"school"`
	Mascot       string   `json:"mascot"`
	Abbreviation string   `json:"abbreviation"`
	Conference   string   `json:"conference"`
	Division     string   `json:"division"`
	Color        string   `json:"color"`
	AltColor     string   `json:"alt_color"`
	Logos        []string `json:"logos"`
}

func RetrieveTeams(year string) {
	url := TeamsURL + "?year=" + year
	fmt.Printf("Retrieving teams data from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var schools []School
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//bodyString := string(bodyBytes)
		fmt.Println(bodyBytes)

		json.Unmarshal(bodyBytes, &schools)

		fmt.Printf("Schools : %+v \n", schools)

		// get sport level
		sportlevel, err := dbsportlevel.ReadByID(dbsportlevel.COLLEGE_FOOTBALL)
		if err != nil {
			log.Fatalf("Error getting sport level record : %v", err)
		}

		fmt.Printf("SportLevel : %#v", sportlevel)

		// get current season record
		season, err := dbseason.ReadCurrentSeason(sportlevel.CurrentYear, sportlevel.SportLevelID)
		if err != nil {
			log.Fatalf("Error getting current season : %v", err)
		}

		for _, school := range schools {
			fmt.Println("=========")
			fmt.Printf("School : %#v \n", school)

			// get the conference id from short_name
			conference, err := dbconference.ReadByExternalNameAndSportLevel(school.Conference, sportlevel.SportLevelID)
			if err != nil {
				log.Fatalf("ERROR: Conference %s not found! : %v", school.Conference, err)
			}

			// get conference_season record or create new
			conferenceSeason, err := dbconferenceseason.ReadByConfSeason(conference.ConferenceID, season.SeasonID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Fatalf("ERROR: %v", err)
				}
				// create a new
				conferenceSeason = &dbconferenceseason.ConferenceSeason{
					ConferenceID: conference.ConferenceID,
					SeasonID:     season.SeasonID,
					DisplayOrder: database.NullInt64{},
				}

				err := dbconferenceseason.Insert(conferenceSeason)
				if err != nil {
					log.Fatalf("ERROR: Can't create ConferenceSeason %v", err)
				}
			}

			// get the team id or create new
			team, err := dbteam.ReadByAbbreviationAndMascot(school.Abbreviation, school.Mascot)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Fatalf("ERROR: %v", err)
				}
				//fmt.Println("No Team with that info")
				log.Fatalf("No Team with that info")
			}

			// get the team_conference_season record or create new
			_, err = dbteamconferenceseason.ReadByTeamConfSeasonID(team.TeamID, conference.ConferenceID, season.SeasonID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Fatalf("ERROR: %v", err)
				}
				// create new
				teamConfSeason := &dbteamconferenceseason.TeamConferenceSeason{
					TeamID:       team.TeamID,
					ConferenceID: conference.ConferenceID,
					SeasonID:     season.SeasonID,
					DivisionName: database.ToNullString(school.Division, true),
				}
				err := dbteamconferenceseason.Insert(teamConfSeason)
				if err != nil {
					log.Fatalf("ERROR: Can't create TeamConferenceSeason %v", err)
				}
			}

		}
	}

}
