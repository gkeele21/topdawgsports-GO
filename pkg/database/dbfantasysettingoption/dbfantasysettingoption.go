package dbfantasysettingoption

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type FantasySettingOption struct {
	FantasySettingOptionID int64              `db:"fantasy_setting_option_id"`
	FantasySettingID       int64              `db:"fantasy_setting_id"`
	Value                  string             `db:"value"`
	DisplayOrder           database.NullInt64 `db:"display_order"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasySettingOption, error) {
	d := FantasySettingOption{}
	err := database.Get(&d, "SELECT * FROM fantasy_setting_option where fantasy_setting_option_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasySettingOption, error) {
	var recs []FantasySettingOption
	err := database.Select(&recs, "SELECT * FROM fantasy_setting_option")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasySettingOption) error {
	_, err := database.Exec("DELETE FROM fantasy_setting_option WHERE fantasy_setting_option_id = ?", d.FantasySettingOptionID)
	if err != nil {
		return fmt.Errorf("fantasy_setting_option: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasySettingOption) error {
	res, err := database.Exec(database.BuildInsert("fantasy_setting_option", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_setting_option: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_setting_option: couldn't get last inserted ID %S", err)
	}

	d.FantasySettingOptionID = ID

	return nil
}
