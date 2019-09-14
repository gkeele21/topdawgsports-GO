package dbteamstandings

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type TeamStandings struct {
	TeamStandingsID        int64                `db:"team_standings_id"`
	TeamConferenceSeasonID int64                `db:"team_conference_season_id"`
	WeekNumber             database.NullInt64   `db:"week_number"`
	Wins                   database.NullInt64   `db:"wins"`
	Losses                 database.NullInt64   `db:"losses"`
	Ties                   database.NullInt64   `db:"ties"`
	WinPercentage          database.NullFloat64 `db:"win_percentage"`
	PointsFor              database.NullFloat64 `db:"points_for"`
	PointsAgainst          database.NullFloat64 `db:"points_against"`
	HomeWins               database.NullInt64   `db:"home_wins"`
	HomeLosses             database.NullInt64   `db:"home_losses"`
	AwayWins               database.NullInt64   `db:"away_wins"`
	AwayLosses             database.NullInt64   `db:"away_losses"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*TeamStandings, error) {
	d := TeamStandings{}
	err := database.Get(&d, "SELECT * FROM team_standings where team_standings_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]TeamStandings, error) {
	var recs []TeamStandings
	err := database.Select(&recs, "SELECT * FROM team_standings")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *TeamStandings) error {
	_, err := database.Exec("DELETE FROM team_standings WHERE team_standings_id = ?", d.TeamStandingsID)
	if err != nil {
		return fmt.Errorf("teamstandings: couldn't delete team_standings %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *TeamStandings) error {
	res, err := database.Exec(database.BuildInsert("team_standings", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("teamstandings: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("teamstandings: couldn't get last inserted ID %S", err)
	}

	d.TeamStandingsID = ID

	return nil
}

// Update will update a record in the database
func Update(s *TeamStandings) error {
	sql := database.BuildUpdate("team_standings", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("team_standings: couldn't update %s", err)
	}

	return nil
}

func Save(s *TeamStandings) error {
	if s.TeamStandingsID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
