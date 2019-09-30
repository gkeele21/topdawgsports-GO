package dbloveemstandings

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type LoveEmStandings struct {
	LoveEmStandingsID int64                `db:"loveem_standings_id"`
	FantasyTeamID     int64                `db:"fantasy_team_id"`
	WeekID            int64                `db:"week_id"`
	WeekGamePts       database.NullFloat64 `db:"week_game_pts"`
	TotalGamePts      database.NullFloat64 `db:"total_game_pts"`
	LeagueRanking     database.NullInt64   `db:"league_ranking"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*LoveEmStandings, error) {
	d := LoveEmStandings{}
	err := database.Get(&d, "SELECT * FROM loveem_standings where loveem_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]LoveEmStandings, error) {
	var recs []LoveEmStandings
	err := database.Select(&recs, "SELECT * FROM loveem_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *LoveEmStandings) error {
	_, err := database.Exec("DELETE FROM loveem_standings WHERE loveem_standings_id = ?", d.LoveEmStandingsID)
	if err != nil {
		return fmt.Errorf("loveem_standings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *LoveEmStandings) error {
	res, err := database.Exec(database.BuildInsert("loveem_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("loveem_standings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("loveem_standings: couldn't get last inserted ID %S", err)
	}

	d.LoveEmStandingsID = ID

	return nil
}

// Update will update a record in the database
func Update(s *LoveEmStandings) error {
	sql := database.BuildUpdate("loveem_standings", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("loveem_standings: couldn't update %s", err)
	}

	return nil
}

func Save(s *LoveEmStandings) error {
	if s.LoveEmStandingsID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadLeagueWeekStandings reads all records in the database for a given leagueId and weekId
func ReadLeagueWeekStandings(leagueID, weekID int64) ([]LoveEmStandings, error) {
	var recs []LoveEmStandings
	err := database.Select(&recs, "SELECT * FROM loveem_standings WHERE fantasy_league_id = ? AND week_id = ?", leagueID, weekID)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
