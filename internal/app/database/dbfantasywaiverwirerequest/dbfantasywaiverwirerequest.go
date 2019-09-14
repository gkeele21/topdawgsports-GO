package dbfantasywaiverwirerequest

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"time"
)

type FantasyWaiverWireRequest struct {
	FantasyWaiverWireRequestID int64               `db:"fantasy_waiver_wire_request_id"`
	FantasyTeamID              int64               `db:"fantasy_team_id"`
	WeekID                     int64               `db:"week_id"`
	Rank                       int64               `db:"rank"`
	RequestDate                time.Time           `db:"request_date"`
	DropPlayerID               database.NullInt64  `db:"drop_player_id"`
	DropType                   database.NullString `db:"drop_type"`
	PickupPlayerID             database.NullInt64  `db:"pu_player_id"`
	PickupType                 database.NullString `db:"pu_type"`
	Processed                  int64               `db:"processed"`
	Granted                    int64               `db:"granted"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyWaiverWireRequest, error) {
	d := FantasyWaiverWireRequest{}
	err := database.Get(&d, "SELECT * FROM fantasy_waiver_wire_request where fantasy_waiver_wire_request_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyWaiverWireRequest, error) {
	var recs []FantasyWaiverWireRequest
	err := database.Select(&recs, "SELECT * FROM fantasy_waiver_wire_request")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyWaiverWireRequest) error {
	_, err := database.Exec("DELETE FROM fantasy_waiver_wire_request WHERE fantasy_waiver_wire_request_id = ?", d.FantasyWaiverWireRequestID)
	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_request: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyWaiverWireRequest) error {
	res, err := database.Exec(database.BuildInsert("fantasy_waiver_wire_request", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_request: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_request: couldn't get last inserted ID %S", err)
	}

	d.FantasyWaiverWireRequestID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyWaiverWireRequest) error {
	sql := database.BuildUpdate("fantasy_waiver_wire_request", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_waiver_wire_request: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyWaiverWireRequest) error {
	if s.FantasyWaiverWireRequestID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
