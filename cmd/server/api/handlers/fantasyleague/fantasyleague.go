package fantasyleague

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbfantasyleague"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbfantasyteam"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type LeagueData struct {
	FantasyLeagueID int64
	SeasonID        int64
	FantasyGameID   int64
	Name            string
	Description     string
	LeaguePassword  string
	Visibility      string
	Status          string
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(g *echo.Group) {
	g.GET("/fantasyleagues/:fantasyLeagueId", getFantasyLeagueByID)
	g.GET("/fantasyleagues/:fantasyLeagueId/teams", getFantasyTeams)
	g.PUT("/fantasyleagues/:fantasyLeagueId", saveLeagueByID)
	g.POST("/fantasyleagues", createLeague)
}

// getFantasyLeagueByID searches for a single fantasy league by leagueid from the route parameter :fantasyLeagueId
func getFantasyLeagueByID(req echo.Context) error {
	var err error

	log.LogRequestData(req)
	searchID := req.Param("fantasyLeagueId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad fantasy league ID given")
	}

	includeVerboseData := req.QueryParam("includeVerboseData")

	var l *dbfantasyleague.FantasyLeague
	if includeVerboseData == "true" {

		//l, err = dbseason.ReadByIDWithSportLevel(num)
		//
		//if err != nil {
		//	return req.Error(http.StatusInternalServerError, "couldn't get season", err)
		//}

	} else {
		l, err = dbfantasyleague.ReadByID(num)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get fantasy league", err)
		}

	}
	return req.JSON(http.StatusOK, l)
}

// getFantasyTeams grabs all fantasy_teams for the given fantasyLeagueId
func getFantasyTeams(req echo.Context) error {
	log.LogRequestData(req)
	tempLeagueID := req.Param("fantasyLeagueId")
	leagueID, err := strconv.ParseInt(tempLeagueID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "pass in a valid integer for fantasyLeagueId", err)
	}

	orderBy := req.QueryParam("orderBy")
	orderByAsc := req.QueryParam("orderByAsc")
	switch orderByAsc {
	case "asc", "ASC", "desc", "DESC":
	default:
		orderByAsc = ""
	}
	// make sure we only allow certain values to be passed in to order by
	switch orderBy {
	case "fantasy_team_id", "user_id", "name", "created_date", "schedule_team_number", "status":
		orderBy += " " + orderByAsc
	default:
		orderBy = ""
	}

	teams, err := dbfantasyteam.ReadAllByFantasyLeagueID(leagueID, orderBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find teams with the fantasy league id.", err)
	}

	return req.JSON(http.StatusOK, teams)

}

// saveLeagueByID searches for a single fantasy_league by fantasyLeagueId from the route parameter :fantasyLeagueId and saves it with the data passed in
func saveLeagueByID(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	searchID := req.Param("fantasyLeagueId")
	fmt.Printf("fantasyLeagueId passed in : %s\n", searchID)
	leagueId, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		req.Logger().Errorf("Error getting leagueId passed in : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "bad fantasyLeagueId given")
	}

	league, err := dbfantasyleague.ReadByID(leagueId)

	if err != nil {
		req.Logger().Errorf("Error getting fantasyleague from db : %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get fantasy_league", err)
	}

	tempLeague := new(LeagueData)
	if err = req.Bind(tempLeague); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}
	if err != nil {
		req.Logger().Errorf("Error populating tempLeague struct : %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), 400)
	}

	fmt.Printf("TempLeague : %#v\n", tempLeague)

	if tempLeague.Name != "" {
		league.Name = tempLeague.Name
	}
	if tempLeague.Description != "" {
		league.Description = database.ToNullString(tempLeague.Description, true)
	}

	league.LeaguePassword = database.ToNullString(tempLeague.LeaguePassword, true)
	league.FantasyGameID = tempLeague.FantasyGameID
	league.SeasonID = tempLeague.SeasonID
	league.Status = tempLeague.Status
	league.Visibility = tempLeague.Visibility

	ret := dbfantasyleague.Update(league)
	if ret != nil {
		req.Logger().Errorf("Error updating fantasyleague record : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, league)
}

// createLeague creates a new fantasyleague with the data passed in
func createLeague(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	league := new(dbfantasyleague.FantasyLeague)

	tempLeague := new(LeagueData)
	if err = req.Bind(tempLeague); err != nil {
		fmt.Printf("Error : %#v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}

	fmt.Printf("TempLeague : %#v\n", tempLeague)

	league.Name = tempLeague.Name
	league.Description = database.ToNullString(tempLeague.Description, true)
	league.LeaguePassword = database.ToNullString(tempLeague.LeaguePassword, true)
	league.FantasyGameID = tempLeague.FantasyGameID
	league.SeasonID = tempLeague.SeasonID
	league.Status = tempLeague.Status
	league.Visibility = tempLeague.Visibility

	ret := dbfantasyleague.Insert(league)
	if ret != nil {
		req.Logger().Errorf("Error updating fantasyleague record : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, league)
}
