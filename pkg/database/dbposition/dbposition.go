package dbposition

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type Position struct {
	PositionID int64               `db:"position_id"`
	SportID    int64               `db:"sport_id"`
	Name       string              `db:"name"`
	NameLong   database.NullString `db:"name_long"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*Position, error) {
	d := Position{}
	err := database.Get(&d, "SELECT * FROM position where position_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]Position, error) {
	var recs []Position
	err := database.Select(&recs, "SELECT * FROM position")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *Position) error {
	_, err := database.Exec("DELETE FROM position WHERE position_id = ?", d.PositionID)
	if err != nil {
		return fmt.Errorf("position: couldn't delete position %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *Position) error {
	res, err := database.Exec(database.BuildInsert("position", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("position: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("position: couldn't get last inserted ID %S", err)
	}

	d.PositionID = ID

	return nil
}
