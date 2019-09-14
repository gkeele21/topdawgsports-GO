package main

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/pkg/stats/football/college/games"
	"github.com/gkeele21/topdawgsportsAPI/internal/pkg/stats/football/college/players/playerstats"
)

func main() {
	fmt.Println("Starting script to import college football player stats...")

	//playerstats.RetrieveGamePlayerStats("2019", "401114223")

	year := "2019"
	week := "02"
	//games.RetrieveGames("2019")

	gameIds := games.RetrieveGamesForWeek(year, week)
	for _, gameId := range gameIds {
		fmt.Printf("Getting stats for gameId %s \n", gameId)
		playerstats.RetrieveGamePlayerStats(year, gameId)
	}

}
