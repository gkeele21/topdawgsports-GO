package dbconferenceseason

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type ConferenceSeason struct {
	ConferenceSeasonID int64              `db:"conference_season_id"`
	ConferenceID       int64              `db:"conference_id"`
	SeasonID           int64              `db:"season_id"`
	DisplayOrder       database.NullInt64 `db:"display_order"`
}

// ReadByID reads conferenceseason by id column
func ReadByID(ID int64) (*ConferenceSeason, error) {
	y := ConferenceSeason{}
	err := database.Get(&y, "SELECT * FROM conference_season where conference_season_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &y, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]ConferenceSeason, error) {
	var recs []ConferenceSeason
	err := database.Select(&recs, "SELECT * FROM conference_season")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a conferenceseason from the database
func Delete(y *ConferenceSeason) error {
	_, err := database.Exec("DELETE FROM conference_season WHERE conference_season_id = ?", y.ConferenceSeasonID)
	if err != nil {
		return fmt.Errorf("conferenceseason: couldn't delete conferenceseason %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(y *ConferenceSeason) error {
	res, err := database.Exec(database.BuildInsert("conference_season", y), database.GetArguments(*y)...)

	if err != nil {
		return fmt.Errorf("conferenceseason: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("conferenceseason: couldn't get last inserted ID %S", err)
	}

	y.ConferenceSeasonID = ID

	return nil
}

// Update will update a record in the database
func Update(s *ConferenceSeason) error {
	sql := database.BuildUpdate("conference_season", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("conference_season: couldn't update %s", err)
	}

	return nil
}

func Save(s *ConferenceSeason) error {
	if s.ConferenceSeasonID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByConfShortName reads conferenceseason by the short_name of the conf and year
func ReadByConfShortName(shortName, year string) (*ConferenceSeason, error) {
	y := ConferenceSeason{}
	err := database.Get(&y, "SELECT cy.* FROM conference_season cy "+
		" INNER JOIN conference c ON c.conference_id = cy.conference_id"+
		" WHERE c.short_name = ?"+
		" AND cy.", shortName)
	if err != nil {
		return nil, err
	}

	return &y, nil
}

// ReadByConfDivSeason reads conferenceseason by the conference_id, season_id columns
func ReadByConfSeason(confId, seasonId int64) (*ConferenceSeason, error) {
	y := ConferenceSeason{}
	err := database.Get(&y, "SELECT * FROM conference_season WHERE conference_id = ? AND season_id = ?", confId, seasonId)
	if err != nil {
		return nil, err
	}

	return &y, nil
}
