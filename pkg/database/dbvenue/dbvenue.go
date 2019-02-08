package dbvenue

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type Venue struct {
	VenueID int64               `db:"venue_id"`
	Name    string              `db:"name"`
	City    database.NullString `db:"city"`
	State   database.NullString `db:"state"`
}

// ReadByID reads venue by id column
func ReadByID(ID int64) (*Venue, error) {
	v := Venue{}
	err := database.Get(&v, "SELECT * FROM venue where venue_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

// ReadAll reads all venues in the database
func ReadAll() ([]Venue, error) {
	var vens []Venue
	err := database.Select(&vens, "SELECT * FROM venue")
	if err != nil {
		return nil, err
	}

	return vens, nil
}

// Delete deletes a venue from the database
func Delete(v *Venue) error {
	_, err := database.Exec("DELETE FROM venue WHERE venue_id = ?", v.VenueID)
	if err != nil {
		return fmt.Errorf("venue: couldn't delete venue %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(v *Venue) error {
	res, err := database.Exec(database.BuildInsert("venue", v), database.GetArguments(*v)...)

	if err != nil {
		return fmt.Errorf("venue: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("venue: couldn't get last inserted ID %S", err)
	}

	v.VenueID = ID

	return nil
}
