package dbseason

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type Season struct {
	SeasonID     int64              `db:"season_id"`
	Name         string             `db:"name"`
	StartingYear database.NullInt64 `db:"starting_year"`
	SportLevelID int64              `db:"sport_level_id"`
	Status       string             `db:"status"`
}

// ReadByID reads season by id column
func ReadByID(ID int64) (*Season, error) {
	s := Season{}
	err := database.Get(&s, "SELECT * FROM season where season_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// ReadAll reads all seasons in the database
func ReadAll() ([]Season, error) {
	var seasons []Season
	err := database.Select(&seasons, "SELECT * FROM season")
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

// Delete deletes a season from the database
func Delete(s *Season) error {
	_, err := database.Exec("DELETE FROM season WHERE season_id = ?", s.SeasonID)
	if err != nil {
		return fmt.Errorf("season: couldn't delete season %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(s *Season) error {
	res, err := database.Exec(database.BuildInsert("season", s), database.GetArguments(*s)...)

	if err != nil {
		return fmt.Errorf("season: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("season: couldn't get last inserted ID %S", err)
	}

	s.SeasonID = ID

	return nil
}
