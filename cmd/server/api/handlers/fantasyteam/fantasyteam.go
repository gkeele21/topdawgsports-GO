package fantasyteam

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbfantasyteam"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type TeamData struct {
	FantasyTeamID      int64
	FantasyLeagueID    int64
	UserID             int64
	Name               string
	Status             string
	ScheduleTeamNumber int64
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(g *echo.Group) {
	g.GET("/fantasyteams/:fantasyTeamId", getFantasyTeamByID)
	g.PUT("/fantasyteams/:fantasyTeamId", saveTeamByID)
	g.POST("/fantasyteams", createTeam)
}

// getFantasyTeamByID searches for a single fantasy team by teamid from the route parameter :fantasyTeamId
func getFantasyTeamByID(req echo.Context) error {
	var err error

	log.LogRequestData(req)
	searchID := req.Param("fantasyTeamId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad fantasy team ID given")
	}

	includeVerboseData := req.QueryParam("includeVerboseData")

	var l *dbfantasyteam.FantasyTeam
	if includeVerboseData == "true" {

		//l, err = dbseason.ReadByIDWithSportLevel(num)
		//
		//if err != nil {
		//	return req.Error(http.StatusInternalServerError, "couldn't get season", err)
		//}

	} else {
		l, err = dbfantasyteam.ReadByID(num)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get fantasy team", err)
		}

	}
	return req.JSON(http.StatusOK, l)
}

// saveTeamByID searches for a single fantasy_team by fantasyTeamId from the route parameter :fantasyTeamId and saves it with the data passed in
func saveTeamByID(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	searchID := req.Param("fantasyTeamId")
	fmt.Printf("fantasyTeamId passed in : %s\n", searchID)
	teamId, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		req.Logger().Errorf("Error getting teamId passed in : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "bad fantasyTeamId given")
	}

	team, err := dbfantasyteam.ReadByID(teamId)

	if err != nil {
		req.Logger().Errorf("Error getting fantasyteam from db : %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get fantasy_team", err)
	}

	tempTeam := new(TeamData)
	if err = req.Bind(tempTeam); err != nil {
		req.Logger().Errorf("Error creating tempTeam from request data : %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}
	if err != nil {
		req.Logger().Errorf("Error populating tempTeam struct : %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), 400)
	}

	fmt.Printf("TempTeam : %#v\n", tempTeam)

	if tempTeam.Name != "" {
		team.Name = tempTeam.Name
	}

	team.FantasyLeagueID = tempTeam.FantasyLeagueID
	team.UserID = tempTeam.UserID
	team.ScheduleTeamNumber = database.ToNullInt(tempTeam.ScheduleTeamNumber, true)
	team.Status = tempTeam.Status

	ret := dbfantasyteam.Update(team)
	if ret != nil {
		req.Logger().Errorf("Error updating fantasyteam record : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, team)
}

// createTeam creates a new fantasyteam with the data passed in
func createTeam(req echo.Context) error {
	var err error
	// Print a copy of this request for debugging.
	log.LogRequestData(req)

	team := new(dbfantasyteam.FantasyTeam)

	tempTeam := new(TeamData)
	if err = req.Bind(tempTeam); err != nil {
		fmt.Printf("Error : %#v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please send a request body", 400)
	}

	fmt.Printf("TempTeam : %#v\n", tempTeam)

	team.Name = tempTeam.Name
	team.FantasyLeagueID = tempTeam.FantasyLeagueID
	team.UserID = tempTeam.UserID
	team.Status = tempTeam.Status
	team.ScheduleTeamNumber = database.ToNullInt(tempTeam.ScheduleTeamNumber, false)
	team.CreatedDate = time.Now()

	fmt.Printf("Team data to create: %#v\n", team)
	ret := dbfantasyteam.Insert(team)
	if ret != nil {
		req.Logger().Errorf("Error updating fantasyteam record : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ret.Error())
	}

	return req.JSON(http.StatusOK, team)
}
