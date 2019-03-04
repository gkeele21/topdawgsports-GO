package main

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handler"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handler/fantasyleague"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handler/season"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handler/sportlevel"
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handler/user"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	// Router
	fmt.Println("API Server now running...")

	port := os.Getenv("TOPDAWG_SERVER_PORT")
	if port == "" {
		port = "8888"
	}
	e := echo.New()
	e.Debug = true

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, host=${host}, uri=${uri}, method=${method}, status=${status}, path=${path}, referer=${referer}, error=${error}, latency=${latency_human},\n",
	}))
	//e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{}))
	e.Use(middleware.Recover())

	registerHandlers(e)

	e.Logger.Fatal(e.Start(":" + port))
}

func registerHandlers(e *echo.Echo) {
	handler.RegisterRoutes(e)
	fantasyleague.RegisterRoutes(e)
	season.RegisterRoutes(e)
	sportlevel.RegisterRoutes(e)
	user.RegisterRoutes(e)
}
