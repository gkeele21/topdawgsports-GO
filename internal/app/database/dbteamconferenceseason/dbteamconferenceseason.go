package dbteamconferenceseason

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type TeamConferenceSeason struct {
	TeamConferenceSeasonID int64               `db:"team_conference_season_id"`
	TeamID                 int64               `db:"team_id"`
	ConferenceID           int64               `db:"conference_id"`
	SeasonID               int64               `db:"season_id"`
	DivisionName           database.NullString `db:"division_name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*TeamConferenceSeason, error) {
	d := TeamConferenceSeason{}
	err := database.Get(&d, "SELECT * FROM team_conference_season where team_conference_season_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]TeamConferenceSeason, error) {
	var recs []TeamConferenceSeason
	err := database.Select(&recs, "SELECT * FROM team_conference_season")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *TeamConferenceSeason) error {
	_, err := database.Exec("DELETE FROM team_conference_season WHERE team_conference_season_id = ?", d.TeamConferenceSeasonID)
	if err != nil {
		return fmt.Errorf("teamconferenceseason: couldn't delete teamconferenceseason %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *TeamConferenceSeason) error {
	res, err := database.Exec(database.BuildInsert("team_conference_season", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("teamconferenceseason: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("teamconferenceseason: couldn't get last inserted ID %S", err)
	}

	d.TeamConferenceSeasonID = ID

	return nil
}

// Update will update a record in the database
func Update(s *TeamConferenceSeason) error {
	sql := database.BuildUpdate("team_conference_season", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("team_conference_season: couldn't update %s", err)
	}

	return nil
}

func Save(s *TeamConferenceSeason) error {
	if s.TeamConferenceSeasonID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByTeamConfSeasonID reads by team_id, conference_id, and season_id columns
func ReadByTeamConfSeasonID(teamId, confId, seasonId int64) (*TeamConferenceSeason, error) {
	d := TeamConferenceSeason{}
	err := database.Get(&d, "SELECT * FROM team_conference_season WHERE team_id = ? AND conference_id = ? AND season_id = ?", teamId, confId, seasonId)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadBySeasonID reads by season_id column
func ReadBySeasonID(seasonId int64) ([]TeamConferenceSeason, error) {
	var recs []TeamConferenceSeason
	err := database.Select(&recs, "SELECT * FROM team_conference_season WHERE season_id = ?", seasonId)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
