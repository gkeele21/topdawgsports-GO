package dbbracketchallengepick

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type BracketChallengePick struct {
	BracketChallengePickID int64              `db:"bracket_challenge_pick_id"`
	FantasyTeamID          int64              `db:"fantasy_team_id"`
	MatchupID              int64              `db:"matchup_id"`
	TeamSeedPickedID       database.NullInt64 `db:"team_seed_picked_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*BracketChallengePick, error) {
	d := BracketChallengePick{}
	err := database.Get(&d, "SELECT * FROM bracket_challenge_pick where bracket_challenge_pick_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]BracketChallengePick, error) {
	var recs []BracketChallengePick
	err := database.Select(&recs, "SELECT * FROM bracket_challenge_pick")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *BracketChallengePick) error {
	_, err := database.Exec("DELETE FROM bracket_challenge_pick WHERE bracket_challenge_pick_id = ?", d.BracketChallengePickID)
	if err != nil {
		return fmt.Errorf("bracketchallengepick: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *BracketChallengePick) error {
	res, err := database.Exec(database.BuildInsert("bracket_challenge_pick", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("bracketchallengepick: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("bracketchallengepick: couldn't get last inserted ID %S", err)
	}

	d.BracketChallengePickID = ID

	return nil
}
