package dbfantasysalarystandings

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasySalaryStandings struct {
	FantasySalaryStandingsID int64              `db:"fantasy_salary_standings_id"`
	FantasyTeamID            int64              `db:"fantasy_team_id"`
	WeekID                   int64              `db:"week_id"`
	WeekFantasyPts           float64            `db:"week_fantasy_pts"`
	TotalFantasyPts          float64            `db:"total_fantasy_pts"`
	WeekGamePts              int64              `db:"week_game_pts"`
	TotalGamePts             int64              `db:"total_game_pts"`
	LeagueRanking            database.NullInt64 `db:"league_ranking"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasySalaryStandings, error) {
	d := FantasySalaryStandings{}
	err := database.Get(&d, "SELECT * FROM fantasy_salary_standings where fantasy_salary_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasySalaryStandings, error) {
	var recs []FantasySalaryStandings
	err := database.Select(&recs, "SELECT * FROM fantasy_salary_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasySalaryStandings) error {
	_, err := database.Exec("DELETE FROM fantasy_salary_standings WHERE fantasy_salary_standings_id = ?", d.FantasySalaryStandingsID)
	if err != nil {
		return fmt.Errorf("fantasy_salary_standings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasySalaryStandings) error {
	res, err := database.Exec(database.BuildInsert("fantasy_salary_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_salary_standings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_salary_standings: couldn't get last inserted ID %S", err)
	}

	d.FantasySalaryStandingsID = ID

	return nil
}
