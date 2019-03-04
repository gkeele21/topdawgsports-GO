package user

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database/dbuser"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type newUserForm struct {
	Email     string
	FirstName string
	LastName  string
	Username  string
	Password  string
	Cell      string
}

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(e *echo.Echo) {
	e.GET("/users/:userId", getUserByID)
	e.GET("/users", getUsers)
	e.POST("/users", newUser)
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

// newUser creates a new user
func newUser(req echo.Context) error {
	//err := req.ParseForm()
	log.LogRequestData(req)

	u := new(newUserForm)
	if err := req.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), 400)
	}

	fmt.Printf("User (populated) : %#v \n", u)

	// Check for required params
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		message := "1 or more required parameters are empty (Email, FirstName, Password)"
		log.Println("Data Input Error", message)
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	if u.Username == "" {
		u.Username = u.Email
	}

	fmt.Printf("UserForm Obj : %#v\n", u)

	// check to see if the username already exists
	exists, err := CheckIfUsernameExists(u.Username)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("ERROR : %#v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error checking for existing username", err)
	}

	if exists {
		fmt.Printf("ERROR : %#v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "submitted username already exists", err)
	}

	user := dbuser.User{
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      database.ToNullString(u.LastName, false),
		Username:      database.ToNullString(u.Username, false),
		UserPassword:  database.ToNullString(u.Password, false),
		Cell:          database.ToNullString(u.Cell, false),
		CreatedDate:   time.Now(),
		LastLoginDate: database.ToNullTime(time.Now(), false),
	}

	fmt.Printf("User Obj : %#v\n", user)

	err = dbuser.Insert(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't insert user", err)
	}

	return req.JSON(http.StatusOK, user)
}

func CheckIfUsernameExists(username string) (bool, error) {
	exists := false
	existingUser, err := dbuser.ReadByUsername(username)
	if err != nil {
		return false, err
	}

	if existingUser.UserID > 0 {
		exists = true
	}

	return exists, nil
}
