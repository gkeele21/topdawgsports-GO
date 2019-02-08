package dbfantasyplayervalue

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type FantasyPlayerValue struct {
	FantasyPlayerValueID int64   `db:"fantasy_player_value_id"`
	FantasyGameID        int64   `db:"fantasy_game_id"`
	PlayerID             int64   `db:"player_id"`
	WeekID               int64   `db:"week_id"`
	Value                float64 `db:"player_value"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyPlayerValue, error) {
	d := FantasyPlayerValue{}
	err := database.Get(&d, "SELECT * FROM fantasy_player_value where fantasy_player_value_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyPlayerValue, error) {
	var recs []FantasyPlayerValue
	err := database.Select(&recs, "SELECT * FROM fantasy_player_value")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyPlayerValue) error {
	_, err := database.Exec("DELETE FROM fantasy_player_value WHERE fantasy_player_value_id = ?", d.FantasyPlayerValueID)
	if err != nil {
		return fmt.Errorf("fantasy_player_value: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyPlayerValue) error {
	res, err := database.Exec(database.BuildInsert("fantasy_player_value", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_player_value: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_player_value: couldn't get last inserted ID %S", err)
	}

	d.FantasyPlayerValueID = ID

	return nil
}
