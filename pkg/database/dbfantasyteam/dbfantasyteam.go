package dbfantasyteam

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
	"time"
)

type FantasyTeam struct {
	FantasyTeamID      int64              `db:"fantasy_team_id"`
	FantasyLeagueID    int64              `db:"fantasy_league_id"`
	UserID             int64              `db:"user_id"`
	Name               string             `db:"name"`
	Created            time.Time          `db:"created_date"`
	Status             string             `db:"status"`
	ScheduleTeamNumber database.NullInt64 `db:"schedule_team_number"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyTeam, error) {
	d := FantasyTeam{}
	err := database.Get(&d, "SELECT * FROM fantasy_team where fantasy_team_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyTeam, error) {
	var recs []FantasyTeam
	err := database.Select(&recs, "SELECT * FROM fantasy_team")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyTeam) error {
	_, err := database.Exec("DELETE FROM fantasy_team WHERE fantasy_team_id = ?", d.FantasyTeamID)
	if err != nil {
		return fmt.Errorf("fantasyteam: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyTeam) error {
	res, err := database.Exec(database.BuildInsert("fantasy_team", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyteam: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyteam: couldn't get last inserted ID %S", err)
	}

	d.FantasyTeamID = ID

	return nil
}

// ReadAllByFantasyLeagueID reads all fantasy_teams in the database for the given fantasyLeagueID
func ReadAllByFantasyLeagueID(fantasyLeagueID int64, orderBy string) ([]FantasyTeam, error) {
	var recs []FantasyTeam
	if orderBy == "" {
		orderBy = "fantasy_team_id asc"
	}
	err := database.Select(&recs, "SELECT * FROM fantasy_team WHERE fantasy_league_id = ? ORDER BY "+orderBy, fantasyLeagueID)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
