package main

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/pkg/stats/football/college/players"
)

func main() {
	fmt.Println("Starting script to import college football players...")

	//players.RetrievePlayers("TCU")
	players.RetrieveAllPlayers()
}
