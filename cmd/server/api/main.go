package main

import (
	"fmt"
	"net/http"
	"topdawgsportsAPI/cmd/server/api/handler"
	"topdawgsportsAPI/cmd/server/api/handler/user"
	"topdawgsportsAPI/cmd/server/api/middleware"

	"topdawgsportsAPI/cmd/server/api/handler/season"

	"github.com/MordFustang21/nova"
	"topdawgsportsAPI/cmd/server/api/handler/sportlevel"
)

func main() {
	// Router
	fmt.Println("API Server now running...")
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
	sportlevel.RegisterRoutes(s)
}
