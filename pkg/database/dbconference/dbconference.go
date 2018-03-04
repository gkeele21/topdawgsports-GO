package dbconference

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type Conference struct {
	ConferenceID int64               `db:"conference_id"`
	Name         string              `db:"name"`
	Abbreviation database.NullString `db:"abbreviation"`
}

// ReadByID reads conference by id column
func ReadByID(ID int64) (*Conference, error) {
	c := Conference{}
	err := database.Get(&c, "SELECT * FROM conference where conference_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// ReadAll reads all conferences in the database
func ReadAll() ([]Conference, error) {
	var confs []Conference
	err := database.Select(&confs, "SELECT * FROM conference")
	if err != nil {
		return nil, err
	}

	return confs, nil
}

// Delete deletes a conference from the database
func Delete(c *Conference) error {
	_, err := database.Exec("DELETE FROM conference WHERE conference_id = ?", c.ConferenceID)
	if err != nil {
		return fmt.Errorf("conference: couldn't delete conference %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(c *Conference) error {
	res, err := database.Exec(database.BuildInsert("conference", c), database.GetArguments(*c)...)

	if err != nil {
		return fmt.Errorf("conference: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("conference: couldn't get last inserted ID %S", err)
	}

	c.ConferenceID = ID

	return nil
}
