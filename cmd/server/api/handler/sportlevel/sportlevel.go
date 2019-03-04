package sportlevel

import (
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbsportlevel"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(e *echo.Echo) {
	e.GET("/sportlevels/:sportLevelId", getSportLevelByID)
	e.GET("/sportlevels", getSportLevels)
}

// getSportLevelByID searches for a single sportlevel by sportlevelid from the route parameter :sportLevelId
func getSportLevelByID(req echo.Context) error {
	var err error

	log.LogRequestData(req)
	searchID := req.Param("sportLevelId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad sportLevelID given")
	}

	var s *dbsportlevel.SportLevelFull
	s, err = dbsportlevel.ReadByIDFull(num)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get sportlevel", err)
	}

	return req.JSON(http.StatusOK, s)
}

// getSportLevels grabs all sportlevels
func getSportLevels(req echo.Context) error {
	log.LogRequestData(req)
	orderBy := req.QueryParam("orderBy")
	orderByAsc := req.QueryParam("orderByAsc")

	switch orderByAsc {
	case "asc", "ASC", "desc", "DESC":
	default:
		orderByAsc = ""
	}
	// make sure we only allow certain values to be passed in to order by
	switch orderBy {
	case "sport_level_id", "level", "sport_id":
		orderBy += " " + orderByAsc
	default:
		orderBy = ""
	}

	sportlevels, err := dbsportlevel.ReadAllFull(orderBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find sportlevels", err)
	}

	return req.JSON(http.StatusOK, sportlevels)
}
