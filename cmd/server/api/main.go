package main

import (
	"net/http"
	"topdawgsportsAPI/cmd/server/api/handler"
	"topdawgsportsAPI/cmd/server/api/handler/user"
	"topdawgsportsAPI/cmd/server/api/middleware"

	"github.com/MordFustang21/nova"
	"topdawgsportsAPI/cmd/server/api/handler/season"
)

func main() {
	// Router
	s := nova.New()

	s.Error(func(req *nova.Request, err error) {

	})

	// Middleware
	middleware.Register(s)
	registerHandlers(s)

	// Serve
	http.ListenAndServe(":8888", s)

}

func registerHandlers(s *nova.Server) {
	handler.RegisterRoutes(s)
	user.RegisterRoutes(s)
	season.RegisterRoutes(s)
}
