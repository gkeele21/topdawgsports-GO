package middleware

import (
	"github.com/MordFustang21/nova"
)

// Register builds a context and trace that follow the request as well as
// declares middleware & routes to handle creating a new UserAuth record,
// validating tokens, and generating new tokens
func Register(s *nova.Server) {

	// anything that needs to happen on every route goes here
	s.Use(func(req *nova.Request, next func()) {
		req.Header().Set("Access-Control-Allow-Origin", "*")
		req.Header().Set("Content-Type", "application/json")
		next()
	})
}
