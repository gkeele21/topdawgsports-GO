package authentication

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbrole"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbuser"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbuserrole"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type loginForm struct {
	Username string
	Password string
}

type returnObject struct {
	Token string
	User  returnUser
}

type newUserForm struct {
	Email     string
	FirstName string
	LastName  string
	Username  string
	Password  string
	Cell      string
}

type returnUser struct {
	IsAdmin   bool   `json:"is_admin"`
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Cell      string `json:"cell"`
}

type JwtClaims struct {
	Name string
	jwt.StandardClaims
}

const JWTSecret = "topdawg5521"

// RegisterRoutes sets up routes on a given nova.Server instance
func RegisterRoutes(e *echo.Echo) {
	e.POST("/login", login)
	e.POST("/register", register)
}

func login(req echo.Context) error {
	log.LogRequestData(req)

	u := new(loginForm)
	if err := req.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), 400)
	}

	username := u.Username
	password := u.Password

	if username == "" || password == "" {
		log.Println("ERROR", "Empty username/password supplied")
		return req.String(http.StatusUnauthorized, "Invalid username/password supplied")
	}
	// check creds from db
	user, err := dbuser.ReadByUsername(username)

	if err == nil && user != nil {
		// check password
		userpass := database.ToNullString(password, false)
		if user.UserPassword == userpass {
			token, err := createJwtToken(strconv.FormatInt(user.PersonID, 10), user.FirstName)
			if err != nil {
				log.Println("ERROR", "Error creating JWT Token", err.Error())
				return req.String(http.StatusInternalServerError, "something went wrong")
			}

			returnUser := returnUser{}
			returnUser.UserID = user.PersonID
			returnUser.FirstName = user.FirstName
			returnUser.LastName = user.LastName.String
			returnUser.Email = user.Email
			returnUser.Cell = user.Cell.String
			returnUser.Username = user.Username.String

			// get user roles, if any
			returnUser.IsAdmin = false
			roles, _ := dbuserrole.ReadByUserID(user.PersonID)

			for _, role := range roles {
				roleId := role.RoleID
				dbrole, err := dbrole.ReadByID(roleId)
				if err == nil {
					if dbrole.Name == "superadmin" {
						returnUser.IsAdmin = true
					}
				}
			}

			returnObj := returnObject{}
			returnObj.User = returnUser
			returnObj.Token = token

			log.Println("INFO", "Returning %#v ", returnObj)
			return req.JSON(http.StatusOK, returnObj)
		}
	}

	return req.String(http.StatusUnauthorized, "Invalid username/password supplied")
}

func createJwtToken(userId, name string) (string, error) {
	claims := JwtClaims{
		name,
		jwt.StandardClaims{
			Id:        userId,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// register creates a new user
func register(req echo.Context) error {
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

	token, err := createJwtToken(strconv.FormatInt(user.PersonID, 10), user.FirstName)
	if err != nil {
		log.Println("ERROR", "Error creating JWT Token", err.Error())
		return req.String(http.StatusInternalServerError, "something went wrong")
	}

	returnUser := returnUser{}
	returnUser.UserID = user.PersonID
	returnUser.FirstName = user.FirstName
	returnUser.LastName = user.LastName.String
	returnUser.Email = user.Email
	returnUser.Cell = user.Cell.String
	returnUser.Username = user.Username.String
	returnUser.IsAdmin = true

	returnObj := returnObject{}
	returnObj.User = returnUser
	returnObj.Token = token

	log.Println("INFO", "Returning %#v ", returnObj)
	return req.JSON(http.StatusOK, returnObj)
}

func CheckIfUsernameExists(username string) (bool, error) {
	exists := false
	existingUser, err := dbuser.ReadByUsername(username)
	if err != nil {
		return false, err
	}

	if existingUser.PersonID > 0 {
		exists = true
	}

	return exists, nil
}
