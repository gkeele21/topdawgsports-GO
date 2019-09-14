package dbfantasysetting

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasySetting struct {
	FantasySettingID int64               `db:"fantasy_setting_id"`
	Name             string              `db:"name"`
	DisplayType      database.NullString `db:"display_type"`
	DisplayOrder     database.NullInt64  `db:"display_order"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasySetting, error) {
	d := FantasySetting{}
	err := database.Get(&d, "SELECT * FROM fantasy_setting where fantasy_setting_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasySetting, error) {
	var recs []FantasySetting
	err := database.Select(&recs, "SELECT * FROM fantasy_setting")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasySetting) error {
	_, err := database.Exec("DELETE FROM fantasy_setting WHERE fantasy_setting_id = ?", d.FantasySettingID)
	if err != nil {
		return fmt.Errorf("fantasy_setting: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasySetting) error {
	res, err := database.Exec(database.BuildInsert("fantasy_setting", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_setting: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_setting: couldn't get last inserted ID %S", err)
	}

	d.FantasySettingID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasySetting) error {
	sql := database.BuildUpdate("fantasy_setting", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_setting: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasySetting) error {
	if s.FantasySettingID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
