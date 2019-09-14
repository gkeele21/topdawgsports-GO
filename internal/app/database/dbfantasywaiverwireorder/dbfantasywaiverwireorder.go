package dbfantasywaiverwireorder

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyWaiverWireOrder struct {
	FantasyWaiverWireOrderID int64              `db:"fantasy_waiver_wire_order_id"`
	FantasyLeagueID          int64              `db:"fantasy_league_id"`
	WeekID                   int64              `db:"week_id"`
	OrderNumber              int64              `db:"order_number"`
	FantasyTeamID            database.NullInt64 `db:"fantasy_team_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyWaiverWireOrder, error) {
	d := FantasyWaiverWireOrder{}
	err := database.Get(&d, "SELECT * FROM fantasy_waiver_wire_order where fantasy_waiver_wire_order_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyWaiverWireOrder, error) {
	var recs []FantasyWaiverWireOrder
	err := database.Select(&recs, "SELECT * FROM fantasy_waiver_wire_order")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyWaiverWireOrder) error {
	_, err := database.Exec("DELETE FROM fantasy_waiver_wire_order WHERE fantasy_waiver_wire_order_id = ?", d.FantasyWaiverWireOrderID)
	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_order: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyWaiverWireOrder) error {
	res, err := database.Exec(database.BuildInsert("fantasy_waiver_wire_order", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_order: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_order: couldn't get last inserted ID %S", err)
	}

	d.FantasyWaiverWireOrderID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyWaiverWireOrder) error {
	sql := database.BuildUpdate("fantasy_waiver_wire_order", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_order: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyWaiverWireOrder) error {
	if s.FantasyWaiverWireOrderID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
