package dbconference

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type Conference struct {
	ConferenceID int64               `db:"conference_id"`
	Name         string              `db:"name"`
	Abbreviation database.NullString `db:"abbreviation"`
	ShortName    database.NullString `db:"short_name"`
	ExternalId   database.NullInt64  `db:"external_id"`
	ExternalName database.NullString `db:"external_name"`
	SportLevelId database.NullInt64  `db:"sport_level_id"`
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

// Update will update a record in the database
func Update(s *Conference) error {
	sql := database.BuildUpdate("conference", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("conference: couldn't update %s", err)
	}

	return nil
}

func Save(s *Conference) error {
	if s.ConferenceID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByExternalNameAndSportLevel reads conference by external_name and sport_level_id columns
func ReadByExternalNameAndSportLevel(externalName string, sportLevelId int64) (*Conference, error) {
	c := Conference{}
	err := database.Get(&c, "SELECT * FROM conference WHERE external_name = ? AND sport_level_id = ?", externalName, sportLevelId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
