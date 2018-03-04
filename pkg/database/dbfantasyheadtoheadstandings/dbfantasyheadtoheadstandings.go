package dbfantasyheadtoheadstandings

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyHeadToHeadStandings struct {
	FantasyHeadToHeadStandingsID int64                `db:"fantasy_headtohead_standings_id"`
	FantasyTeamID                int64                `db:"fantasy_team_id"`
	WeekID                       int64                `db:"week_id"`
	WeekFantasyPts               database.NullFloat64 `db:"week_fantasy_pts"`
	TotalFantasyPts              database.NullFloat64 `db:"total_fantasy_pts"`
	WeekFantasyPtsAgainst        database.NullFloat64 `db:"week_fantasy_pts_against"`
	TotalFantasyPtsAgainst       database.NullFloat64 `db:"total_fantasy_pts_against"`
	Wins                         int64                `db:"wins"`
	Losses                       int64                `db:"losses"`
	Ties                         int64                `db:"ties"`
	WeekHiScore                  database.NullInt64   `db:"week_hi_score"`
	TotalHiScore                 database.NullInt64   `db:"total_hi_score"`
	LeagueRanking                database.NullInt64   `db:"league_ranking"`
	CurrentStreak                database.NullInt64   `db:"current_streak"`
	LastFive                     database.NullString  `db:"last_five"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyHeadToHeadStandings, error) {
	d := FantasyHeadToHeadStandings{}
	err := database.Get(&d, "SELECT * FROM fantasy_headtohead_standings where fantasy_headtohead_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyHeadToHeadStandings, error) {
	var recs []FantasyHeadToHeadStandings
	err := database.Select(&recs, "SELECT * FROM fantasy_headtohead_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyHeadToHeadStandings) error {
	_, err := database.Exec("DELETE FROM fantasy_headtohead_standings WHERE fantasy_headtohead_standings_id = ?", d.FantasyHeadToHeadStandingsID)
	if err != nil {
		return fmt.Errorf("fantasy_headtohead_standings: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyHeadToHeadStandings) error {
	res, err := database.Exec(database.BuildInsert("fantasy_headtohead_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_headtohead_standings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_headtohead_standings: couldn't get last inserted ID %S", err)
	}

	d.FantasyHeadToHeadStandingsID = ID

	return nil
}
