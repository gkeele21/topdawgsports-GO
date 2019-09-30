package user

import (
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbfantasyteam"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbuser"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(g *echo.Group) {
	g.GET("/users/:userId", getUserByID)
	g.GET("/users", getUsers)
	g.GET("/users/:userId/activeteams", getActiveTeams)
}

// Response is the json representation of a user
//type Response struct {
//	User dbuser.User
//}

// getUserByID searches for a single user by user id from the route parameter :userId
func getUserByID(req echo.Context) error {
	var err error

	log.LogRequestData(req)
	searchID := req.Param("userId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad user ID given")
	}

	var u *dbuser.User
	u, err = dbuser.ReadByID(num)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get user", err)
	}

	return req.JSON(http.StatusOK, u)
}

// getUsers grabs all users
func getUsers(req echo.Context) error {
	log.LogRequestData(req)
	users, err := dbuser.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find users", err)
	}

	return req.JSON(http.StatusOK, users)
}

// getActiveTeams grabs all active teams for the user
func getActiveTeams(req echo.Context) error {
	log.LogRequestData(req)
	userID := req.Param("userId")
	num, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad user ID given")
	}

	var u *dbuser.User
	u, err = dbuser.ReadByID(num)

	users, err := dbfantasyteam.ReadByUserIDFull(u.UserID, "active","")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find users", err)
	}

	return req.JSON(http.StatusOK, users)
}
