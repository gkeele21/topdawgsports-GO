package user

import (
	"github.com/MordFustang21/nova"
	"net/http"
	"strconv"
	"topdawgsportsAPI/pkg/database/dbuser"
	"topdawgsportsAPI/pkg/log"
	"topdawgsportsAPI/pkg/database"
	"time"
	"fmt"
	"database/sql"
)

type newUserForm struct {
	Email     string
	FirstName string
	LastName  string
	Username  string
	Password  string
	Cell      string
}

// RegisterRoutes sets up routs on a given nova.Server instance
func RegisterRoutes(s *nova.Server) {
	s.Get("/user/:userId", getUserByID)
	s.Get("/user/list", getUsers)
	s.Post("/user", newUser)
}

// Response is the json representation of a user
//type Response struct {
//	User dbuser.User
//}

// getUserByID searches for a single user by user id from the route parameter :userId
func getUserByID(req *nova.Request) error {
	var err error

	log.LogRequest(req)
	searchID := req.RouteParam("userId")
	num, err := strconv.ParseInt(searchID, 10, 64)
	if err != nil {
		return req.Error(http.StatusBadRequest, "bad user ID given")
	}

	var u *dbuser.User
	u, err = dbuser.ReadByID(num)

	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't get user", err)
	}

	return req.JSON(http.StatusOK, u)
}

// getUsers grabs all users
func getUsers(req *nova.Request) error {
	log.LogRequest(req)
	users, err := dbuser.ReadAll()
	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't find users", err)
	}

	return req.JSON(http.StatusOK, users)
}

// newUser creates a new user
func newUser(req *nova.Request) error {
	log.LogRequest(req)

	u := newUserForm{}

	err := req.ReadJSON(&u)
	if err != nil {
		fmt.Printf("ERROR : %#v\n", err)
		return req.Error(http.StatusBadRequest, "invalid json body provided for login form", err)
	}

	log.LogRequestData(u)

	// Check for required params
	if u.FirstName == "" || u.Email == "" || u.Password == "" {
		message := "1 or more required parameters are empty (Email, FirstName, Password)"
		log.Println("Data Input Error", message)
		return req.Error(http.StatusBadRequest, message)
	}

	if u.Username == "" {
		u.Username = u.Email
	}

	fmt.Printf("UserForm Obj : %#v\n", u)

	// check to see if the username already exists
	exists, err := CheckIfUsernameExists(u.Username)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("ERROR : %#v\n", err)
		return req.Error(http.StatusInternalServerError, "error checking for existing username", err)
	}

	if exists {
		fmt.Printf("ERROR : %#v\n", err)
		return req.Error(http.StatusBadRequest, "submitted username already exists", err)
	}

	user := dbuser.User{
		Email:        database.ToNullString(u.Email, false),
		FirstName:    database.ToNullString(u.FirstName, false),
		LastName:     database.ToNullString(u.LastName, false),
		Username:     database.ToNullString(u.Username, false),
		UserPassword: database.ToNullString(u.Password, false),
		Cell:         database.ToNullString(u.Cell, false),
		Created:      time.Now(),
		LastLogin:    time.Now(),
	}

	fmt.Printf("User Obj : %#v\n", user)

	err = dbuser.Insert(&user)
	if err != nil {
		return req.Error(http.StatusInternalServerError, "couldn't insert user", err)
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
