package handler

import (
	"net/http"

	"github.com/MordFustang21/nova"
)

// RegisterRoutes adds routes to provided mux
func RegisterRoutes(s *nova.Server) {
	s.Get("/health", healthCheck)
}

func healthCheck(req *nova.Request) error {
	//req.StatusCode(http.StatusOK)

	return req.JSON(http.StatusOK, "ok")
	//return req.Send("ok")
}
