package dbmatchup

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type Matchup struct {
	MatchupID     int64              `db:"matchup_id"`
	WeekID        database.NullInt64 `db:"week_id"`
	MatchupDate   database.NullTime  `db:"matchup_date"`
	VisitorTeamID database.NullInt64 `db:"visitor_team_id"`
	HomeTeamID    database.NullInt64 `db:"home_team_id"`
	VenueID       database.NullInt64 `db:"venue_id"`
	VisitorScore  database.NullInt64 `db:"visitor_score"`
	HomeScore     database.NullInt64 `db:"home_score"`
	WinningTeamID database.NullInt64 `db:"winning_team_id"`
	NumOvertimes  database.NullInt64 `db:"num_overtimes"`
	Status        string             `db:"status"`
	ExternalId    database.NullInt64 `db:"external_id"`
}

const STATUS_PENDING = "pending"
const STATUS_ACTIVE = "active"
const STATUS_FINAL = "final"

// ReadByID reads matchup by id column
func ReadByID(ID int64) (*Matchup, error) {
	m := Matchup{}
	err := database.Get(&m, "SELECT * FROM matchup where matchup_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// ReadAll reads all matchups in the database
func ReadAll() ([]Matchup, error) {
	var ms []Matchup
	err := database.Select(&ms, "SELECT * FROM matchup")
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// Delete deletes a matchup from the database
func Delete(m *Matchup) error {
	_, err := database.Exec("DELETE FROM matchup WHERE matchup_id = ?", m.MatchupID)
	if err != nil {
		return fmt.Errorf("matchup: couldn't delete matchup %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(m *Matchup) error {
	res, err := database.Exec(database.BuildInsert("matchup", m), database.GetArguments(*m)...)

	if err != nil {
		return fmt.Errorf("matchup: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("matchup: couldn't get last inserted ID %S", err)
	}

	m.MatchupID = ID

	return nil
}

// Update will update a record in the database
func Update(s *Matchup) error {
	sql := database.BuildUpdate("matchup", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("matchup: couldn't update %s", err)
	}

	return nil
}

func Save(s *Matchup) error {
	if s.MatchupID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByWeekIDAndVisitorIDAndHomeId reads matchup by week_id, visitor_id, and home_id columns
func ReadByWeekIDAndVisitorIDAndHomeId(weekId, visitorId, homeId int64) (*Matchup, error) {
	m := Matchup{}
	err := database.Get(&m, "SELECT * FROM matchup WHERE week_id = ? AND visitor_team_id = ? AND home_team_id = ?", weekId, visitorId, homeId)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// ReadByExternalID reads matchup by external_id column
func ReadByExternalID(ID int64) (*Matchup, error) {
	m := Matchup{}
	err := database.Get(&m, "SELECT * FROM matchup where external_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (matchup *Matchup) SetWinner() {
	if matchup.HomeScore.Valid && matchup.VisitorScore.Valid {
		if matchup.HomeScore.Int64 > matchup.VisitorScore.Int64 {
			(*matchup).WinningTeamID = matchup.HomeTeamID
			(*matchup).Status = STATUS_FINAL
		} else if matchup.VisitorScore.Int64 > matchup.HomeScore.Int64 {
			(*matchup).WinningTeamID = matchup.VisitorTeamID
			(*matchup).Status = STATUS_FINAL
		}
	}

}
