package dbpickempick

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type PickEmPick struct {
	PickEmPickID     int64              `db:"pickem_pick_id"`
	FantasyTeamID    int64              `db:"fantasy_team_id"`
	WeekID           int64              `db:"week_id"`
	MatchupID        int64              `db:"matchup_id"`
	TeamPickedID     int64              `db:"team_picked_id"`
	ConfidencePoints database.NullInt64 `db:"confidence_points"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*PickEmPick, error) {
	d := PickEmPick{}
	err := database.Get(&d, "SELECT * FROM pickem_pick where pickem_pick_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]PickEmPick, error) {
	var recs []PickEmPick
	err := database.Select(&recs, "SELECT * FROM pickem_pick")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *PickEmPick) error {
	_, err := database.Exec("DELETE FROM pickem_pick WHERE pickem_pick_id = ?", d.PickEmPickID)
	if err != nil {
		return fmt.Errorf("pickem_pick: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *PickEmPick) error {
	res, err := database.Exec(database.BuildInsert("pickem_pick", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("pickem_pick: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("pickem_pick: couldn't get last inserted ID %S", err)
	}

	d.PickEmPickID = ID

	return nil
}
