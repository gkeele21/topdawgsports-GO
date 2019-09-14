package dbloveempickteam

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type LoveEmPickTeam struct {
	LoveEmPickTeamID int64 `db:"loveem_pick_team_id"`
	FantasyTeamID    int64 `db:"fantasy_team_id"`
	WeekID           int64 `db:"week_id"`
	TeamPickedID     int64 `db:"team_picked_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*LoveEmPickTeam, error) {
	d := LoveEmPickTeam{}
	err := database.Get(&d, "SELECT * FROM loveem_pick_team where loveem_pick_team_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]LoveEmPickTeam, error) {
	var recs []LoveEmPickTeam
	err := database.Select(&recs, "SELECT * FROM loveem_pick_team")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *LoveEmPickTeam) error {
	_, err := database.Exec("DELETE FROM loveem_pick_team WHERE loveem_pick_team_id = ?", d.LoveEmPickTeamID)
	if err != nil {
		return fmt.Errorf("loveem_pick_team: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *LoveEmPickTeam) error {
	res, err := database.Exec(database.BuildInsert("loveem_pick_team", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("loveem_pick_team: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("loveem_pick_team: couldn't get last inserted ID %S", err)
	}

	d.LoveEmPickTeamID = ID

	return nil
}

// Update will update a record in the database
func Update(s *LoveEmPickTeam) error {
	sql := database.BuildUpdate("loveem_pick_team", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("loveem_pick_team: couldn't update %s", err)
	}

	return nil
}

func Save(s *LoveEmPickTeam) error {
	if s.LoveEmPickTeamID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
