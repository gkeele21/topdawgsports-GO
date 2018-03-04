package dbfantasyroster

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyRoster struct {
	FantasyRosterID   int64              `db:"fantasy_roster_id"`
	FantasyTeamID     int64              `db:"fantasy_team_id"`
	WeekID            int64              `db:"week_id"`
	PlayerID          database.NullInt64 `db:"player_id"`
	ScoringState      string             `db:"scoring_state"`
	ScoringPositionID database.NullInt64 `db:"scoring_position_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyRoster, error) {
	d := FantasyRoster{}
	err := database.Get(&d, "SELECT * FROM fantasy_roster where fantasy_roster_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyRoster, error) {
	var recs []FantasyRoster
	err := database.Select(&recs, "SELECT * FROM fantasy_roster")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyRoster) error {
	_, err := database.Exec("DELETE FROM fantasy_roster WHERE fantasy_roster_id = ?", d.FantasyRosterID)
	if err != nil {
		return fmt.Errorf("fantasyroster: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyRoster) error {
	res, err := database.Exec(database.BuildInsert("fantasy_roster", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyroster: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyroster: couldn't get last inserted ID %S", err)
	}

	d.FantasyRosterID = ID

	return nil
}
