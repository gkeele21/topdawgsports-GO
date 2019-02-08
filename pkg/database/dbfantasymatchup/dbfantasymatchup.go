package dbfantasymatchup

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type FantasyMatchup struct {
	FantasyMatchupID int64                `db:"fantasy_matchup_id"`
	FantasyLeagueID  int64                `db:"fantasy_league_id"`
	WeekID           int64                `db:"week_id"`
	VisitorTeamID    database.NullInt64   `db:"visitor_team_id"`
	HomeTeamID       database.NullInt64   `db:"home_team_id"`
	VisitorScore     database.NullFloat64 `db:"visitor_score"`
	HomeScore        database.NullFloat64 `db:"home_score"`
	WinningTeamID    database.NullInt64   `db:"winning_team_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyMatchup, error) {
	d := FantasyMatchup{}
	err := database.Get(&d, "SELECT * FROM fantasy_matchup where fantasy_matchup_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyMatchup, error) {
	var recs []FantasyMatchup
	err := database.Select(&recs, "SELECT * FROM fantasy_matchup")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyMatchup) error {
	_, err := database.Exec("DELETE FROM fantasy_matchup WHERE fantasy_matchup_id = ?", d.FantasyMatchupID)
	if err != nil {
		return fmt.Errorf("fantasymatchup: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyMatchup) error {
	res, err := database.Exec(database.BuildInsert("fantasy_matchup", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasymatchup: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasymatchup: couldn't get last inserted ID %S", err)
	}

	d.FantasyMatchupID = ID

	return nil
}
