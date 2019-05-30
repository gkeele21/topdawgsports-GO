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
	CreatedDate        time.Time          `db:"created_date"`
	Status             string             `db:"status"`
	ScheduleTeamNumber database.NullInt64 `db:"schedule_team_number"`
}

type FantasyTeamFull struct {
	FantasyTeamID      int64              `db:"fantasy_team_id"`
	TeamName           string             `db:"fantasy_team_name"`
	DateCreated        time.Time          `db:"created_date"`
	Status             string             `db:"status"`
	ScheduleTeamNumber database.NullInt64 `db:"schedule_team_number"`
	FantasyLeagueID    int64              `db:"fantasy_league_id"`
	FantasyLeagueName  string             `db:"fantasy_league_name"`
	UserID             int64              `db:"user_id"`
	UserFirstName      string             `db:"first_name"`
	UserLastName       string             `db:"last_name"`
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

// Update will update a record in the database
func Update(s *FantasyTeam) error {
	sql := database.BuildUpdate("fantasy_team", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("season: couldn't update %s", err)
	}

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

// ReadAllByFantasyLeagueIDFull reads all fantasy_teams in the database for the given fantasyLeagueID
func ReadAllByFantasyLeagueIDFull(fantasyLeagueID int64, orderBy string) ([]FantasyTeamFull, error) {
	var recs []FantasyTeamFull
	if orderBy == "" {
		orderBy = "fantasy_team_id asc"
	}
	sql := "SELECT ft.fantasy_team_id, ft.name as fantasy_team_name, ft.created_date, ft.status, ft.schedule_team_number, ft.fantasy_league_id, fl.name as fantasy_league_name, ft.user_id, u.first_name, u.last_name" +
		" FROM fantasy_team ft " +
		" INNER JOIN user u ON u.user_id = ft.user_id " +
		" INNER JOIN fantasy_league fl ON fl.fantasy_league_id = ft.fantasy_league_id " +
		" WHERE ft.fantasy_league_id = ? ORDER BY " + orderBy
	fmt.Printf("SQL : %s\n", sql)
	fmt.Printf("FantasyLeagueID: %s\n", fantasyLeagueID)
	err := database.Select(&recs, sql, fantasyLeagueID)
	if err != nil {
		fmt.Printf("Error getting teams: %#v\n", err)
		return nil, err
	}

	return recs, nil
}
