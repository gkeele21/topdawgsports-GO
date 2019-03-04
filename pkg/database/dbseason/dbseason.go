package dbseason

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
)

type Season struct {
	SeasonID     int64              `db:"season_id"`
	Name         string             `db:"name"`
	StartingYear database.NullInt64 `db:"starting_year"`
	SportLevelID int64              `db:"sport_level_id"`
	Status       string             `db:"status"`
}

type SeasonSportLevel struct {
	SeasonID     int64              `db:"season_id"`
	Name         string             `db:"name"`
	StartingYear database.NullInt64 `db:"starting_year"`
	SportLevelID int64              `db:"sport_level_id"`
	Status       string             `db:"status"`
	SportID      int64              `db:"sport_id"`
	SportLevel   string             `db:"level"`
	SportName    string             `db:"sport_name"`
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

// ReadByIDWithSportLevel reads season by id column
func ReadByIDWithSportLevel(ID int64) (*SeasonSportLevel, error) {
	s := SeasonSportLevel{}
	err := database.Get(&s, "SELECT se.*, sl.level, s.sport_id, s.name as 'sport_name' FROM season se INNER JOIN sport_level sl ON sl.sport_level_id = se.sport_level_id INNER JOIN sport s ON s.sport_id = sl.sport_id WHERE se.season_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// ReadAll reads all seasons in the database
func ReadAll(orderBy string) ([]Season, error) {
	var seasons []Season
	if orderBy == "" {
		orderBy = "season_id asc"
	}
	err := database.Select(&seasons, "SELECT * FROM season ORDER BY "+orderBy)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

// ReadAll reads all seasons in the database and the matching sportlevel and sport data
func ReadAllWithSportLevel(orderBy string) ([]SeasonSportLevel, error) {
	var seasons []SeasonSportLevel
	if orderBy == "" {
		orderBy = "season_id asc"
	}
	err := database.Select(&seasons, "SELECT se.*, sl.level, s.sport_id, s.name as 'sport_name' FROM season se INNER JOIN sport_level sl ON sl.sport_level_id = se.sport_level_id INNER JOIN sport s ON s.sport_id = sl.sport_id ORDER BY "+orderBy)
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

// Update will update a record in the database
func Update(s *Season) error {
	sql := database.BuildUpdate("season", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("season: couldn't update %s", err)
	}

	return nil
}
