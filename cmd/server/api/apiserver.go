package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/authentication"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/basic"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/fantasyleague"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/season"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/sportlevel"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/user"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/middlewares"
	"github.com/gkeele21/topdawgsportsAPI/pkg/log"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

func main() {
	// Router
	log.Println("INFO", "API Server now running...")

	port := os.Getenv("TOPDAWG_SERVER_PORT")
	if port == "" {
		port = "8888"
	}

	e := New()

	e.Logger.Fatal(e.Start(":" + port))
}

func New() *echo.Echo {
	e := echo.New()
	e.Debug = true

	// Create Groups
	jwtGroup := e.Group("/api/v1")

	// Middleware
	middlewares.SetMainMiddleware(e)
	middlewares.SetJwtMiddleware(jwtGroup)

	// Handlers
	jwtGroup.GET("/test", testJwt)

	registerHandlers(e, jwtGroup)

	return e
}

func testJwt(req echo.Context) error {
	log.Println("INFO", "in mainJwt")
	user := req.Get("user")
	log.Println("INFO", "User obj %#v", user)

	jwtToken := user.(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)

	fmt.Println("User name : ", claims["Name"], "User id : ", claims["jti"])

	return req.String(http.StatusOK, "You are on the main JWT page")

}

func registerHandlers(e *echo.Echo, g *echo.Group) {
	basic.RegisterRoutes(e)
	authentication.RegisterRoutes(e)
	fantasyleague.RegisterRoutes(g)
	season.RegisterRoutes(g)
	sportlevel.RegisterRoutes(g)
	user.RegisterRoutes(g)
}

