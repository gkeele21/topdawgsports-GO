package dbfantasygamesetting

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyGameSetting struct {
	FantasyGameSettingID int64              `db:"fantasy_game_setting_id"`
	FantasyGameID        int64              `db:"fantasy_game_id"`
	PositionID           database.NullInt64 `db:"position_id"`
	FantasySettingID     int64              `db:"fantasy_setting_id"`
	Value                string             `db:"value"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyGameSetting, error) {
	d := FantasyGameSetting{}
	err := database.Get(&d, "SELECT * FROM fantasy_game_setting where fantasy_game_setting_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyGameSetting, error) {
	var recs []FantasyGameSetting
	err := database.Select(&recs, "SELECT * FROM fantasy_game_setting")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyGameSetting) error {
	_, err := database.Exec("DELETE FROM fantasy_game_setting WHERE fantasy_game_setting_id = ?", d.FantasyGameSettingID)
	if err != nil {
		return fmt.Errorf("fantasy_game_setting: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyGameSetting) error {
	res, err := database.Exec(database.BuildInsert("fantasy_game_setting", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_game_setting: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_game_setting: couldn't get last inserted ID %S", err)
	}

	d.FantasyGameSettingID = ID

	return nil
}
