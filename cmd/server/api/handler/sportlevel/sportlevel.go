package sportlevel

import (
	"fmt"
	"github.com/MordFustang21/nova"
	"net/http"
	"strconv"
	"topdawgsportsAPI/pkg/database/dbsportlevel"
	"topdawgsportsAPI/pkg/log"
)

// RegisterRoutes sets up routs on a given nova.Server instance
func RegisterRoutes(s *nova.Server) {
	s.Get("/sportlevels/:sportLevelId", getSportLevelByID)
	s.Get("/sportlevels", getSportLevels)
}

// getSportLevelByID searches for a single sportlevel by sportlevelid from the route parameter :sportLevelId
func getSportLevelByID(req *nova.Request) error {
	var err error

	log.LogRequest(req)
	searchID := req.RouteParam("sportLevelId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad sportLevelID given")
	}

	var s *dbsportlevel.SportLevelFull
	s, err = dbsportlevel.ReadByIDFull(num)

	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't get sportlevel", err)
	}

	return req.JSON(http.StatusOK, s)
}

// getSportLevels grabs all sportlevels
func getSportLevels(req *nova.Request) error {
	log.LogRequest(req)
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

	fmt.Printf("OrderBy %s ", orderBy)
	sportlevels, err := dbsportlevel.ReadAllFull(orderBy)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't find sportlevels", err)
	}

	return req.JSON(http.StatusOK, sportlevels)
}
