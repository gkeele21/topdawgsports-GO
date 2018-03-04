package dbpickemstandings

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type PickEmStandings struct {
	PickEmStandingsID int64              `db:"pickem_standings_id"`
	FantasyTeamID     int64              `db:"fantasy_team_id"`
	WeekID            int64              `db:"week_id"`
	WeekGamePoints    int64              `db:"week_game_pts"`
	TotalGamePoints   int64              `db:"total_game_pts"`
	LeagueRanking     database.NullInt64 `db:"league_ranking"`
	GamesCorrect      database.NullInt64 `db:"games_correct"`
	TotalGamesCorrect database.NullInt64 `db:"total_games_correct"`
	GamesWrong        database.NullInt64 `db:"games_wrong"`
	TotalGamesWrong   database.NullInt64 `db:"total_games_wrong"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*PickEmStandings, error) {
	d := PickEmStandings{}
	err := database.Get(&d, "SELECT * FROM pickem_standings where pickem_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]PickEmStandings, error) {
	var recs []PickEmStandings
	err := database.Select(&recs, "SELECT * FROM pickem_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *PickEmStandings) error {
	_, err := database.Exec("DELETE FROM pickem_standings WHERE pickem_standings_id = ?", d.PickEmStandingsID)
	if err != nil {
		return fmt.Errorf("pickem_standings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *PickEmStandings) error {
	res, err := database.Exec(database.BuildInsert("pickem_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("pickem_standings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("pickem_standings: couldn't get last inserted ID %S", err)
	}

	d.PickEmStandingsID = ID

	return nil
}
