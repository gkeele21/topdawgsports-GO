package fantasyleague

import (
	"encoding/json"
	"fmt"
	"github.com/MordFustang21/nova"
	"net/http"
	"strconv"
	"topdawgsportsAPI/pkg/database/dbfantasyleague"
	"topdawgsportsAPI/pkg/database/dbfantasyteam"
	"topdawgsportsAPI/pkg/log"
)

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(s *nova.Server) {
	s.Get("/fantasyleagues/:fantasyLeagueId", getFantasyLeagueByID)
	s.Get("/fantasyleagues/:fantasyLeagueId/teams", getFantasyTeams)
	s.Post("/fantasyleagues/:fantasyLeagueId", saveLeagueByID)
}

// getFantasyLeagueByID searches for a single fantasy league by leagueid from the route parameter :fantasyLeagueId
func getFantasyLeagueByID(req *nova.Request) error {
	var err error

	log.LogRequestData(req)
	searchID := req.RouteParam("fantasyLeagueId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad fantasy league ID given")
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
			return req.Error(http.StatusInternalServerError, "couldn't get fantasy league", err)
		}

	}
	return req.JSON(http.StatusOK, l)
}

// getFantasyTeams grabs all fantasy_teams for the given fantasyLeagueId
func getFantasyTeams(req *nova.Request) error {
	log.LogRequestData(req)
	tempLeagueID := req.RouteParam("fantasyLeagueId")
	leagueID, err := strconv.ParseInt(tempLeagueID, 10, 64)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "pass in a valid integer for fantasyLeagueId", err)
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
		return req.Error(http.StatusInternalServerError, "couldn't find teams with the fantasy league id.", err)
	}

	return req.JSON(http.StatusOK, teams)

}

// saveLeagueByID searches for a single fantasy_league by fantasyLeagueId from the route parameter :fantasyLeagueId and saves it with the data passed in
func saveLeagueByID(req *nova.Request) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	searchID := req.RouteParam("fantasyLeagueId")
	fmt.Printf("fantasyLeagueId passed in : %s\n", searchID)
	leagueId, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad fantasyLeagueId given")
	}

	_, err = dbfantasyleague.ReadByID(leagueId)

	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't get fantasy_league", err)
	}

	if req.Body == nil {
		return req.Error(http.StatusInternalServerError, "Please send a request body", 400)
	}

	var tempLeague *dbfantasyleague.FantasyLeague
	err = json.NewDecoder(req.Body).Decode(&tempLeague)
	if err != nil {
		return req.Error(http.StatusInternalServerError, err.Error(), 400)
	}

	fmt.Printf("TempLeague : %#v\n", tempLeague)

	ret := dbfantasyleague.Update(tempLeague)
	if ret != nil {
		return req.Error(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, tempLeague)
}


