package main

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/pkg/stats/football/college/teams"
)

func main() {
	fmt.Println("Starting script to import college football teams...")

	teams.RetrieveTeams("2019")
}
