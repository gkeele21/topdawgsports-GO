package dbfantasydraft

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyDraft struct {
	FantasyDraftID  int64              `db:"fantasy_draft_id"`
	FantasyLeagueID int64              `db:"fantasy_league_id"`
	Round           int64              `db:"round"`
	Place           int64              `db:"place"`
	FantasyTeamID   database.NullInt64 `db:"fantasy_team_id"`
	PlayerID        database.NullInt64 `db:"player_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyDraft, error) {
	d := FantasyDraft{}
	err := database.Get(&d, "SELECT * FROM fantasy_draft where fantasy_draft_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyDraft, error) {
	var recs []FantasyDraft
	err := database.Select(&recs, "SELECT * FROM fantasy_draft")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyDraft) error {
	_, err := database.Exec("DELETE FROM fantasy_draft WHERE fantasy_draft_id = ?", d.FantasyDraftID)
	if err != nil {
		return fmt.Errorf("fantasy_draft: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyDraft) error {
	res, err := database.Exec(database.BuildInsert("fantasy_draft", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_draft: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_draft: couldn't get last inserted ID %S", err)
	}

	d.FantasyDraftID = ID

	return nil
}
