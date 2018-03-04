package dbfantasyleaguesetting

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyLeagueSetting struct {
	FantasyLeagueSettingID int64              `db:"fantasy_league_setting_id"`
	FantasyLeagueID        int64              `db:"fantasy_league_id"`
	PositionID             database.NullInt64 `db:"position_id"`
	FantasySettingID       int64              `db:"fantasy_setting_id"`
	Value                  string             `db:"value"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyLeagueSetting, error) {
	d := FantasyLeagueSetting{}
	err := database.Get(&d, "SELECT * FROM fantasy_league_setting where fantasy_league_setting_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyLeagueSetting, error) {
	var recs []FantasyLeagueSetting
	err := database.Select(&recs, "SELECT * FROM fantasy_league_setting")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyLeagueSetting) error {
	_, err := database.Exec("DELETE FROM fantasy_league_setting WHERE fantasy_league_setting_id = ?", d.FantasyLeagueSettingID)
	if err != nil {
		return fmt.Errorf("fantasy_league_setting: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyLeagueSetting) error {
	res, err := database.Exec(database.BuildInsert("fantasy_league_setting", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_league_setting: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_league_setting: couldn't get last inserted ID %S", err)
	}

	d.FantasyLeagueSettingID = ID

	return nil
}
