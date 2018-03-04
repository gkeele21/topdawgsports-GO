package dbsportlevel

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type SportLevel struct {
	SportLevelID int64  `db:"sport_level_id"`
	SportID      int64  `db:"sport_id"`
	Level        string `db:"level"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*SportLevel, error) {
	d := SportLevel{}
	err := database.Get(&d, "SELECT * FROM sport_level where sport_level_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]SportLevel, error) {
	var recs []SportLevel
	err := database.Select(&recs, "SELECT * FROM sport_level")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *SportLevel) error {
	_, err := database.Exec("DELETE FROM sport_level WHERE sport_level_id = ?", d.SportLevelID)
	if err != nil {
		return fmt.Errorf("sport_level: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *SportLevel) error {
	res, err := database.Exec(database.BuildInsert("sport_level", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("sport_level: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("sport_level: couldn't get last inserted ID %S", err)
	}

	d.SportLevelID = ID

	return nil
}
