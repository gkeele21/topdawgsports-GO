package dbbracketchallengestandings

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type BracketChallengeStandings struct {
	BracketChallengeStandingsID int64              `db:"bracket_challenge_standings_id"`
	FantasyTeamID               int64              `db:"fantasy_team_id"`
	WeekID                      int64              `db:"week_id"`
	Rank                        database.NullInt64 `db:"rank"`
	RoundPoints                 int64              `db:"round_points"`
	TotalPoints                 int64              `db:"total_points"`
	MaxPossible                 database.NullInt64 `db:"max_possible"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*BracketChallengeStandings, error) {
	d := BracketChallengeStandings{}
	err := database.Get(&d, "SELECT * FROM bracket_challenge_standings where bracket_challenge_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]BracketChallengeStandings, error) {
	var recs []BracketChallengeStandings
	err := database.Select(&recs, "SELECT * FROM bracket_challenge_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *BracketChallengeStandings) error {
	_, err := database.Exec("DELETE FROM bracket_challenge_standings WHERE bracket_challenge_standings_id = ?", d.BracketChallengeStandingsID)
	if err != nil {
		return fmt.Errorf("bracketchallengestandings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *BracketChallengeStandings) error {
	res, err := database.Exec(database.BuildInsert("bracket_challenge_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("bracketchallengestandings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("bracketchallengestandings: couldn't get last inserted ID %S", err)
	}

	d.BracketChallengeStandingsID = ID

	return nil
}
