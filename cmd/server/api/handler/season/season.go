package season

import (
	"fmt"
	"github.com/MordFustang21/nova"
	"net/http"
	"strconv"
	"topdawgsportsAPI/pkg/database"
	"topdawgsportsAPI/pkg/database/dbfantasyleague"
	"topdawgsportsAPI/pkg/database/dbseason"
	"topdawgsportsAPI/pkg/log"
)

type SeasonData struct {
	SeasonId     int64
	Name         string
	Status       string
	SportLevelId int64
	StartingYear int64
}

// RegisterRoutes sets up routs on a given nova.Server instance
func RegisterRoutes(s *nova.Server) {
	s.Get("/seasons/:seasonId", getSeasonByID)
	s.Post("/seasons/:seasonId", saveSeasonByID)
	s.Get("/seasons", getSeasons)
	s.Get("/seasons/:seasonId/games/:gameId/leagues", getSeasonGameLeagues)
}

// getSeasonByID searches for a single season by seasonid from the route parameter :seasonId
func getSeasonByID(req *nova.Request) error {
	var err error

	log.LogRequest(req)
	searchID := req.RouteParam("seasonId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad season ID given")
	}

	includeSportLevels := req.QueryParam("includeSportLevels")

	if includeSportLevels == "true" {
		var s *dbseason.SeasonSportLevel
		s, err = dbseason.ReadByIDWithSportLevel(num)

		if err != nil {
			return req.Error(http.StatusInternalServerError, "couldn't get season", err)
		}

		return req.JSON(http.StatusOK, s)
	} else {
		var s *dbseason.Season
		s, err = dbseason.ReadByID(num)

		if err != nil {
			return req.Error(http.StatusInternalServerError, "couldn't get season", err)
		}

		return req.JSON(http.StatusOK, s)
	}
}

// saveSeasonByID searches for a single season by seasonid from the route parameter :seasonId and saves it with the data passed in
func saveSeasonByID(req *nova.Request) error {
	fmt.Println("In saveSeasonById")
	var err error

	log.LogRequest(req)

	var tempSeason SeasonData
	err = req.ReadJSON(&tempSeason)
	fmt.Printf("TempSeason : %#v\n", tempSeason)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad season data passed in")
	}

	searchID := req.RouteParam("seasonId")
	seasonId, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad season ID given")
	}

	var s *dbseason.Season
	// TODO: after grabbing values from db, update columns that have been passed in to us
	s, err = dbseason.ReadByID(seasonId)

	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't get season", err)
	}

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
		return req.Error(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, s)
}

// getSeasons grabs all seasons
func getSeasons(req *nova.Request) error {
	log.LogRequest(req)
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

	fmt.Printf("OrderBy %s ", orderBy)
	if includeSportLevels == "true" {
		seasons, err := dbseason.ReadAllWithSportLevel(orderBy)
		if err != nil {
			return req.Error(http.StatusInternalServerError, "couldn't find seasons", err)
		}

		return req.JSON(http.StatusOK, seasons)
	} else {
		seasons, err := dbseason.ReadAll(orderBy)
		if err != nil {
			return req.Error(http.StatusInternalServerError, "couldn't find seasons", err)
		}

		return req.JSON(http.StatusOK, seasons)
	}

}

// getSeasonGameLeagues grabs all fantasy_leagues for the given seasonId and gameId
func getSeasonGameLeagues(req *nova.Request) error {
	log.LogRequest(req)
	tempSeasonID := req.RouteParam("seasonId")
	seasonID, err := strconv.ParseInt(tempSeasonID, 10, 64)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "pass in a valid integer for seasonId", err)
	}

	tempGameID := req.RouteParam("gameId")
	gameID, err := strconv.ParseInt(tempGameID, 10, 64)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "pass in a valid integer for fantasyGameId", err)
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

	fmt.Printf("OrderBy %s ", orderBy)
	leagues, err := dbfantasyleague.ReadAllBySeasonIDFantasyGameID(seasonID, gameID, orderBy)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't find leagues with the season id.", err)
	}

	return req.JSON(http.StatusOK, leagues)

}
