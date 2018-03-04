package dbteamconferenceyear

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type TeamConferenceYear struct {
	TeamConferenceYearID int64 `db:"team_conference_year_id"`
	ConferenceID         int64 `db:"conference_year_id"`
	TeamID               int64 `db:"team_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*TeamConferenceYear, error) {
	d := TeamConferenceYear{}
	err := database.Get(&d, "SELECT * FROM team_conference_year where team_conference_year_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]TeamConferenceYear, error) {
	var recs []TeamConferenceYear
	err := database.Select(&recs, "SELECT * FROM team_conference_year")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *TeamConferenceYear) error {
	_, err := database.Exec("DELETE FROM team_conference_year WHERE team_conference_year_id = ?", d.TeamConferenceYearID)
	if err != nil {
		return fmt.Errorf("teamconferenceyear: couldn't delete teamconferenceyear %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *TeamConferenceYear) error {
	res, err := database.Exec(database.BuildInsert("team_conference_year", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("teamconferenceyear: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("teamconferenceyear: couldn't get last inserted ID %S", err)
	}

	d.TeamConferenceYearID = ID

	return nil
}
