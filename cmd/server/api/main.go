package main

import (
	"topdawgsportsAPI/cmd/server/api/handler"
	"topdawgsportsAPI/cmd/server/api/handler/user"
	"topdawgsportsAPI/cmd/server/api/middleware"
	"net/http"

	"github.com/MordFustang21/nova"
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
}
