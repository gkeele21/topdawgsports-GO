package dbfantasyschedulebase

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyScheduleBase struct {
	FantasyScheduleBaseID int64  `db:"fantasy_schedule_base_id"`
	Name                  string `db:"name"`
	WeekName              string `db:"week_name"`
	GameNumber            int64  `db:"game_number"`
	HomeTeam              int64  `db:"home_team"`
	AwayTeam              int64  `db:"away_team"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyScheduleBase, error) {
	d := FantasyScheduleBase{}
	err := database.Get(&d, "SELECT * FROM fantasy_schedule_base where fantasy_schedule_base_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyScheduleBase, error) {
	var recs []FantasyScheduleBase
	err := database.Select(&recs, "SELECT * FROM fantasy_schedule_base")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyScheduleBase) error {
	_, err := database.Exec("DELETE FROM fantasy_schedule_base WHERE fantasy_schedule_base_id = ?", d.FantasyScheduleBaseID)
	if err != nil {
		return fmt.Errorf("fantasy_schedule_base: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyScheduleBase) error {
	res, err := database.Exec(database.BuildInsert("fantasy_schedule_base", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_schedule_base: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_schedule_base: couldn't get last inserted ID %S", err)
	}

	d.FantasyScheduleBaseID = ID

	return nil
}
