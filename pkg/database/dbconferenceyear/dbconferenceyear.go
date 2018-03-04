package dbconferenceyear

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type ConferenceYear struct {
	ConferenceYearID int64              `db:"conference_year_id"`
	ConferenceID     int64              `db:"conference_id"`
	DivisionID       database.NullInt64 `db:"division_id"`
	SeasonID         int64              `db:"season_id"`
	DisplayOrder     database.NullInt64 `db:"display_order"`
}

// ReadByID reads conferenceyear by id column
func ReadByID(ID int64) (*ConferenceYear, error) {
	y := ConferenceYear{}
	err := database.Get(&y, "SELECT * FROM conference_year where conference_year_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &y, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]ConferenceYear, error) {
	var recs []ConferenceYear
	err := database.Select(&recs, "SELECT * FROM conference_year")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a division from the database
func Delete(y *ConferenceYear) error {
	_, err := database.Exec("DELETE FROM conference_year WHERE conference_year_id = ?", y.ConferenceYearID)
	if err != nil {
		return fmt.Errorf("conferenceyear: couldn't delete conferenceyear %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(y *ConferenceYear) error {
	res, err := database.Exec(database.BuildInsert("conference_year", y), database.GetArguments(*y)...)

	if err != nil {
		return fmt.Errorf("conferenceyear: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("conferenceyear: couldn't get last inserted ID %S", err)
	}

	y.ConferenceYearID = ID

	return nil
}
