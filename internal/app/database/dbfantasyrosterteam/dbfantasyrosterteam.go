package dbfantasyrosterteam

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyRosterTeam struct {
	FantasyRosterTeamID int64  `db:"fantasy_roster_team_id"`
	FantasyTeamID       int64  `db:"fantasy_team_id"`
	WeekID              int64  `db:"week_id"`
	TeamID              int64  `db:"team_id"`
	ScoringState        string `db:"scoring_state"`
}

const SCORINGSTATE_SCORING = "scoring"
const SCORINGSTATE_BENCH = "bench"
const SCORINGSTATE_IR = "ir"

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyRosterTeam, error) {
	d := FantasyRosterTeam{}
	err := database.Get(&d, "SELECT * FROM fantasy_roster_team where fantasy_roster_team_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyRosterTeam, error) {
	var recs []FantasyRosterTeam
	err := database.Select(&recs, "SELECT * FROM fantasy_roster_team")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyRosterTeam) error {
	_, err := database.Exec("DELETE FROM fantasy_roster_team WHERE fantasy_roster_team_id = ?", d.FantasyRosterTeamID)
	if err != nil {
		return fmt.Errorf("fantasyrosterteam: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyRosterTeam) error {
	res, err := database.Exec(database.BuildInsert("fantasy_roster_team", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyrosterteam: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyrosterteam: couldn't get last inserted ID %S", err)
	}

	d.FantasyRosterTeamID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyRosterTeam) error {
	sql := database.BuildUpdate("fantasy_roster_team", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_roster_team: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyRosterTeam) error {
	if s.FantasyRosterTeamID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
