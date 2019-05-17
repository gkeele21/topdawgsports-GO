package season

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbfantasyleague"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbseason"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type SeasonData struct {
	SeasonId     int64
	Name         string
	Status       string
	SportLevelId int64
	StartingYear int64
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(g *echo.Group) {
	g.GET("/seasons/:seasonId", getSeasonByID)
	g.PUT("/seasons/:seasonId", saveSeasonByID)
	g.POST("/seasons", createSeason)
	g.GET("/seasons", getSeasons)
	g.GET("/seasons/:seasonId/games/:gameId/leagues", getSeasonGameLeagues)
}

// getSeasonByID searches for a single season by seasonid from the route parameter :seasonId
func getSeasonByID(req echo.Context) error {
	var err error

	fmt.Println("Here in getSeasonByID")
	log.LogRequestData(req)
	searchID := req.Param("seasonId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad season ID given")
	}

	includeSportLevels := req.QueryParam("includeSportLevels")

	if includeSportLevels == "true" {
		var s *dbseason.SeasonSportLevel
		fmt.Println("Calling dbseason.ReadByIDWithSportLevel")
		s, err = dbseason.ReadByIDWithSportLevel(num)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get season", err)
		}

		return req.JSON(http.StatusOK, s)
	} else {
		var s *dbseason.Season
		fmt.Println("Calling dbseason.ReadById")
		s, err = dbseason.ReadByID(num)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get season", err)
		}

		return req.JSON(http.StatusOK, s)
	}
}

// saveSeasonByID searches for a single season by seasonid from the route parameter :seasonId and saves it with the data passed in
func saveSeasonByID(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	searchID := req.Param("seasonId")
	seasonId, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad season ID given")
	}

	var s *dbseason.Season
	s, err = dbseason.ReadByID(seasonId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get season", err)
	}

	tempSeason := new(SeasonData)
	if err = req.Bind(tempSeason); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}

	fmt.Printf("TempSeason : %#v\n", tempSeason)

	if tempSeason.Name != "" {
		s.Name = tempSeason.Name
	}

	s.Status = tempSeason.Status
	if tempSeason.StartingYear > 0 {
		s.StartingYear = database.ToNullInt(tempSeason.StartingYear, false)
	}
	s.SportLevelID = tempSeason.SportLevelId

	fmt.Printf("Season data to update: %#v\n", s)
	ret := dbseason.Update(s)
	if ret != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, s)
}

// getSeasons grabs all seasons
func getSeasons(req echo.Context) error {
	log.LogRequestData(req)
	orderBy := req.QueryParam("orderBy")
	orderByAsc := req.QueryParam("orderByAsc")
	includeSportLevels := req.QueryParam("includeSportLevels")

	switch orderByAsc {
	case "asc", "ASC", "desc", "DESC":
	default:
		orderByAsc = ""
	}
	// make sure we only allow certain values to be passed in to order by
	switch orderBy {
	case "season_id", "name", "starting_year", "sport_level_id", "status":
		orderBy += " " + orderByAsc
	default:
		orderBy = ""
	}

	if includeSportLevels == "true" {
		seasons, err := dbseason.ReadAllWithSportLevel(orderBy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find seasons", err)
		}

		return req.JSON(http.StatusOK, seasons)
	} else {
		seasons, err := dbseason.ReadAll(orderBy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find seasons", err)
		}

		return req.JSON(http.StatusOK, seasons)
	}

}

// getSeasonGameLeagues grabs all fantasy_leagues for the given seasonId and gameId
func getSeasonGameLeagues(req echo.Context) error {
	log.LogRequestData(req)
	tempSeasonID := req.Param("seasonId")
	seasonID, err := strconv.ParseInt(tempSeasonID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "pass in a valid integer for seasonId", err)
	}

	tempGameID := req.Param("gameId")
	gameID, err := strconv.ParseInt(tempGameID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "pass in a valid integer for fantasyGameId", err)
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
	case "fantasy_league_id", "name", "starting_year", "sport_level_id", "status":
		orderBy += " " + orderByAsc
	default:
		orderBy = ""
	}

	leagues, err := dbfantasyleague.ReadAllBySeasonIDFantasyGameID(seasonID, gameID, orderBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find leagues with the season id.", err)
	}

	return req.JSON(http.StatusOK, leagues)

}

// createSeason creates a new season record with the data passed in
func createSeason(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	tempSeason := new(SeasonData)
	if err = req.Bind(tempSeason); err != nil {
		fmt.Printf("Error : %#v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}

	fmt.Printf("TempSeason : %#v\n", tempSeason)

	s := new(dbseason.Season)
	if tempSeason.Name != "" {
		s.Name = tempSeason.Name
	}

	s.Status = tempSeason.Status
	s.StartingYear = database.ToNullInt(tempSeason.StartingYear, false)
	s.SportLevelID = tempSeason.SportLevelId
	s.Status = tempSeason.Status

	fmt.Printf("Season data to create: %#v\n", s)
	ret := dbseason.Insert(s)
	if ret != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, s)
}
