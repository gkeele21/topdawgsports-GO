package middlewares

import (
	"github.com/gkeele21/topdawgsportsAPI/cmd/server/api/handlers/authentication"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetMainMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, host=${host}, uri=${uri}, method=${method}, status=${status}, path=${path}, referer=${referer}, error=${error}, latency=${latency_human},\n",
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
}

func SetJwtMiddleware(g *echo.Group) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(authentication.JWTSecret),
	}))

}
