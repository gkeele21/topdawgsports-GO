package dbweek

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type Week struct {
	WeekID    int64             `db:"week_id"`
	SeasonID  int64             `db:"season_id"`
	WeekName  string            `db:"week_name"`
	StartDate database.NullTime `db:"start_date"`
	EndDate   database.NullTime `db:"end_date"`
	Status    string            `db:"status"`
	WeekType  string            `db:"week_type"`
}

const STATUS_PENDING = "pending"
const STATUS_ACTIVE = "active"
const STATUS_FINAL = "final"

const TYPE_INITIAL = "initial"
const TYPE_MIDDLE = "middle"
const TYPE_FINAL = "final"

// ReadByID reads week by id column
func ReadByID(ID int64) (*Week, error) {
	w := Week{}
	err := database.Get(&w, "SELECT * FROM week where week_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

// ReadAll reads all weeks in the database
func ReadAll() ([]Week, error) {
	var weeks []Week
	err := database.Select(&weeks, "SELECT * FROM week")
	if err != nil {
		return nil, err
	}

	return weeks, nil
}

// Delete deletes a week from the database
func Delete(w *Week) error {
	_, err := database.Exec("DELETE FROM week WHERE week_id = ?", w.WeekID)
	if err != nil {
		return fmt.Errorf("week: couldn't delete week %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(w *Week) error {
	res, err := database.Exec(database.BuildInsert("week", w), database.GetArguments(*w)...)

	if err != nil {
		return fmt.Errorf("week: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("week: couldn't get last inserted ID %S", err)
	}

	w.WeekID = ID

	return nil
}

// Update will update a record in the database
func Update(s *Week) error {
	sql := database.BuildUpdate("week", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("week: couldn't update %s", err)
	}

	return nil
}

func Save(s *Week) error {
	if s.WeekID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByWeekNameAndSeasonID reads week by week_name and season_id columns
func ReadByWeekNameAndSeasonID(weekName string, seasonId int64) (*Week, error) {
	w := Week{}
	err := database.Get(&w, "SELECT * FROM week where week_name = ? AND season_id = ?", weekName, seasonId)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
