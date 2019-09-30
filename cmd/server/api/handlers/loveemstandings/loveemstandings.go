package loveemstandings

import (
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbloveemstandings"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type StandingsData struct {
	FantasyTeamName string
	UserID          int64
	UserName        int64
	WeekGamePoints  string
	TotalGamePoints string
	LeagueRanking   int64
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(g *echo.Group) {
	g.GET("/loveemstandings/:fantasyLeagueId/:weekId", getWeekStandings)
}

// getWeekStandings retrieves standings results for a single fantasy league by leagueid from the route parameter :fantasyLeagueId
func getWeekStandings(req echo.Context) error {
	var err error

	log.LogRequestData(req)
	tempLeagueID := req.Param("fantasyLeagueId")
	leagueID, err := strconv.ParseInt(tempLeagueID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad fantasy league ID given")
	}
	tempWeekID := req.Param("weekId")
	weekID, err := strconv.ParseInt(tempWeekID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad week ID given")
	}

	var s []dbloveemstandings.LoveEmStandings
	s, err = dbloveemstandings.ReadLeagueWeekStandings(leagueID, weekID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get loveem standings", err)
	}
	return req.JSON(http.StatusOK, s)
}
