package dbfantasyleaguepayoutteam

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyLeaguePayoutTeam struct {
	FantasyLeaguePayoutTeamID int64                `db:"fantasy_league_payout_team_id"`
	FantasyLeaguePayoutID     int64                `db:"fantasy_league_payout_id"`
	FantasyTeamID             database.NullInt64   `db:"fantasy_team_id"`
	PayoutAmount              database.NullFloat64 `db:"payout_amount"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyLeaguePayoutTeam, error) {
	d := FantasyLeaguePayoutTeam{}
	err := database.Get(&d, "SELECT * FROM fantasy_league_payout_team where fantasy_league_payout_team_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyLeaguePayoutTeam, error) {
	var recs []FantasyLeaguePayoutTeam
	err := database.Select(&recs, "SELECT * FROM fantasy_league_payout_team")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyLeaguePayoutTeam) error {
	_, err := database.Exec("DELETE FROM fantasy_league_payout_team WHERE fantasy_league_payout_team_id = ?", d.FantasyLeaguePayoutTeamID)
	if err != nil {
		return fmt.Errorf("fantasyleaguepayoutteam: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyLeaguePayoutTeam) error {
	res, err := database.Exec(database.BuildInsert("fantasy_league_payout_team", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyleaguepayoutteam: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyleaguepayoutteam: couldn't get last inserted ID %S", err)
	}

	d.FantasyLeaguePayoutTeamID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyLeaguePayoutTeam) error {
	sql := database.BuildUpdate("fantasy_league_payout_team", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_league_payout_team: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyLeaguePayoutTeam) error {
	if s.FantasyLeaguePayoutTeamID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
