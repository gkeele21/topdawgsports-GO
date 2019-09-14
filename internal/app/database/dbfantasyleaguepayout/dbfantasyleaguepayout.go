package dbfantasyleaguepayout

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyLeaguePayout struct {
	FantasyLeaguePayoutID int64               `db:"fantasy_league_payout_id"`
	FantasyLeagueID       int64               `db:"fantasy_league_id"`
	Name                  string              `db:"name"`
	Prize                 database.NullString `db:"prize"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyLeaguePayout, error) {
	d := FantasyLeaguePayout{}
	err := database.Get(&d, "SELECT * FROM fantasy_league_payout where fantasy_league_payout_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyLeaguePayout, error) {
	var recs []FantasyLeaguePayout
	err := database.Select(&recs, "SELECT * FROM fantasy_league_payout")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyLeaguePayout) error {
	_, err := database.Exec("DELETE FROM fantasy_league_payout WHERE fantasy_league_payout_id = ?", d.FantasyLeaguePayoutID)
	if err != nil {
		return fmt.Errorf("fantasyleaguepayout: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyLeaguePayout) error {
	res, err := database.Exec(database.BuildInsert("fantasy_league_payout", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyleaguepayout: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyleaguepayout: couldn't get last inserted ID %S", err)
	}

	d.FantasyLeaguePayoutID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyLeaguePayout) error {
	sql := database.BuildUpdate("fantasy_league_payout", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_league_payout: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyLeaguePayout) error {
	if s.FantasyLeaguePayoutID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
