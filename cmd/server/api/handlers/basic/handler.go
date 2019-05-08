package basic

import (
	"github.com/labstack/echo"
	"net/http"
)

// RegisterRoutes adds routes to provided mux
func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", healthCheck)
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
