package fantasyleague

import (
	"github.com/MordFustang21/nova"
	"net/http"
	"strconv"
	"topdawgsportsAPI/pkg/database/dbfantasyleague"
	"topdawgsportsAPI/pkg/log"
)

type LeagueData struct {
	SeasonId     int64
	Name         string
	Status       string
	SportLevelId int64
	StartingYear int64
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(s *nova.Server) {
	s.Get("/fantasyleagues/:fantasyLeagueId", getFantasyLeagueByID)
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

