package dbdivision

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type Division struct {
	DivisionID int64  `db:"division_id"`
	Name       string `db:"name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*Division, error) {
	d := Division{}
	err := database.Get(&d, "SELECT * FROM division where division_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]Division, error) {
	var recs []Division
	err := database.Select(&recs, "SELECT * FROM division")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *Division) error {
	_, err := database.Exec("DELETE FROM division WHERE division_id = ?", d.DivisionID)
	if err != nil {
		return fmt.Errorf("division: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *Division) error {
	res, err := database.Exec(database.BuildInsert("division", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("division: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("division: couldn't get last inserted ID %S", err)
	}

	d.DivisionID = ID

	return nil
}
