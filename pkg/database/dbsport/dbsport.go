package dbsport

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type Sport struct {
	SportID int64  `db:"sport_id"`
	Name    string `db:"name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*Sport, error) {
	d := Sport{}
	err := database.Get(&d, "SELECT * FROM sport where sport_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]Sport, error) {
	var recs []Sport
	err := database.Select(&recs, "SELECT * FROM sport")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *Sport) error {
	_, err := database.Exec("DELETE FROM sport WHERE sport_id = ?", d.SportID)
	if err != nil {
		return fmt.Errorf("sport: couldn't delete sport %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *Sport) error {
	res, err := database.Exec(database.BuildInsert("sport", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("sport: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("sport: couldn't get last inserted ID %S", err)
	}

	d.SportID = ID

	return nil
}
