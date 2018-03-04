package dbfantasyfootballplayerstats

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyFootballPlayerStats struct {
	FantasyFootballPlayerStatsID int64   `db:"fantasy_football_player_stats_id"`
	WeekID                       int64   `db:"week_id"`
	FantasyLeagueID              int64   `db:"fantasy_league_id"`
	PlayerID                     int64   `db:"player_id"`
	FantasyPoints                float64 `db:"fantasy_points"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyFootballPlayerStats, error) {
	d := FantasyFootballPlayerStats{}
	err := database.Get(&d, "SELECT * FROM fantasy_football_player_stats where fantasy_football_player_stats_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyFootballPlayerStats, error) {
	var recs []FantasyFootballPlayerStats
	err := database.Select(&recs, "SELECT * FROM fantasy_football_player_stats")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyFootballPlayerStats) error {
	_, err := database.Exec("DELETE FROM fantasy_football_player_stats WHERE fantasy_football_player_stats_id = ?", d.FantasyFootballPlayerStatsID)
	if err != nil {
		return fmt.Errorf("fantasy_football_player_stats: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyFootballPlayerStats) error {
	res, err := database.Exec(database.BuildInsert("fantasy_football_player_stats", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_football_player_stats: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_football_player_stats: couldn't get last inserted ID %S", err)
	}

	d.FantasyFootballPlayerStatsID = ID

	return nil
}
