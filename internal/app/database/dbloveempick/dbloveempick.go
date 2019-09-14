package dbloveempick

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type LoveEmPick struct {
	LoveEmPickID int64 `db:"loveem_pick_id"`
	FantasyTeamID    int64 `db:"fantasy_team_id"`
	WeekID           int64 `db:"week_id"`
	PlayerPickedID     int64 `db:"player_picked_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*LoveEmPick, error) {
	d := LoveEmPick{}
	err := database.Get(&d, "SELECT * FROM loveem_pick where loveem_pick_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]LoveEmPick, error) {
	var recs []LoveEmPick
	err := database.Select(&recs, "SELECT * FROM loveem_pick")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *LoveEmPick) error {
	_, err := database.Exec("DELETE FROM loveem_pick WHERE loveem_pick_id = ?", d.LoveEmPickID)
	if err != nil {
		return fmt.Errorf("loveem_pick: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *LoveEmPick) error {
	res, err := database.Exec(database.BuildInsert("loveem_pick", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("loveem_pick: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("loveem_pick: couldn't get last inserted ID %S", err)
	}

	d.LoveEmPickID = ID

	return nil
}

// Update will update a record in the database
func Update(s *LoveEmPick) error {
	sql := database.BuildUpdate("loveem_pick", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("loveem_pick: couldn't update %s", err)
	}

	return nil
}

func Save(s *LoveEmPick) error {
	if s.LoveEmPickID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
