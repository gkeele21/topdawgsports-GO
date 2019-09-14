package dbfantasytransaction

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"time"
)

type FantasyTransaction struct {
	FantasyTransactionID int64               `db:"fantasy_transaction_id"`
	FantasyLeagueID      int64               `db:"fantasy_league_id"`
	FantasyTeamID        int64               `db:"fantasy_team_id"`
	WeekID               database.NullInt64  `db:"week_id"`
	TransactionDate      time.Time           `db:"transaction_date"`
	DropType             database.NullString `db:"drop_type"`
	DropPlayerID         database.NullInt64  `db:"drop_player_id"`
	PickupType           database.NullString `db:"pu_type"`
	PickupPlayerID       database.NullInt64  `db:"pu_player_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyTransaction, error) {
	d := FantasyTransaction{}
	err := database.Get(&d, "SELECT * FROM fantasy_transaction where fantasy_transaction_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyTransaction, error) {
	var recs []FantasyTransaction
	err := database.Select(&recs, "SELECT * FROM fantasy_transaction")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyTransaction) error {
	_, err := database.Exec("DELETE FROM fantasy_transaction WHERE fantasy_transaction_id = ?", d.FantasyTransactionID)
	if err != nil {
		return fmt.Errorf("fantasy_transaction: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyTransaction) error {
	res, err := database.Exec(database.BuildInsert("fantasy_transaction", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_transaction: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_transaction: couldn't get last inserted ID %S", err)
	}

	d.FantasyTransactionID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyTransaction) error {
	sql := database.BuildUpdate("fantasy_transaction", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_transaction: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyTransaction) error {
	if s.FantasyTransactionID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
