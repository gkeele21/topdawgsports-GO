package dbfantasysettinggroup

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type FantasySettingGroup struct {
	FantasySettingGroupID int64              `db:"fantasy_setting_group_id"`
	Name                  string             `db:"name"`
	DisplayOrder          database.NullInt64 `db:"display_order"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasySettingGroup, error) {
	d := FantasySettingGroup{}
	err := database.Get(&d, "SELECT * FROM fantasy_setting_group where fantasy_setting_group_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasySettingGroup, error) {
	var recs []FantasySettingGroup
	err := database.Select(&recs, "SELECT * FROM fantasy_setting_group")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasySettingGroup) error {
	_, err := database.Exec("DELETE FROM fantasy_setting_group WHERE fantasy_setting_group_id = ?", d.FantasySettingGroupID)
	if err != nil {
		return fmt.Errorf("fantasy_setting_group: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasySettingGroup) error {
	res, err := database.Exec(database.BuildInsert("fantasy_setting_group", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_setting_group: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_setting_group: couldn't get last inserted ID %S", err)
	}

	d.FantasySettingGroupID = ID

	return nil
}
