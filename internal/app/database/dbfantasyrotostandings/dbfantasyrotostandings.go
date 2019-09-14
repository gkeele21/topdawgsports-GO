package dbfantasyrotostandings

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyRotoStandings struct {
	FantasyRotoStandingsID int64                `db:"fantasy_roto_standings_id"`
	FantasyTeamID          int64                `db:"fantasy_team_id"`
	WeekID                 int64                `db:"week_id"`
	WeekFantasyPts         database.NullFloat64 `db:"week_fantasy_pts"`
	TotalFantasyPts        database.NullFloat64 `db:"total_fantasy_pts"`
	WeekGamePts            database.NullInt64   `db:"week_game_pts"`
	TotalGamePts           database.NullInt64   `db:"total_game_pts"`
	LeagueRanking          database.NullInt64   `db:"league_ranking"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyRotoStandings, error) {
	d := FantasyRotoStandings{}
	err := database.Get(&d, "SELECT * FROM fantasy_roto_standings where fantasy_roto_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyRotoStandings, error) {
	var recs []FantasyRotoStandings
	err := database.Select(&recs, "SELECT * FROM fantasy_roto_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyRotoStandings) error {
	_, err := database.Exec("DELETE FROM fantasy_roto_standings WHERE fantasy_roto_standings_id = ?", d.FantasyRotoStandingsID)
	if err != nil {
		return fmt.Errorf("fantasy_roto_standings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyRotoStandings) error {
	res, err := database.Exec(database.BuildInsert("fantasy_roto_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_roto_standings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_roto_standings: couldn't get last inserted ID %S", err)
	}

	d.FantasyRotoStandingsID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyRotoStandings) error {
	sql := database.BuildUpdate("fantasy_roto_standings", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_roto_standings: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyRotoStandings) error {
	if s.FantasyRotoStandingsID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
