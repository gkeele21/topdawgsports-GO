package dbrankingsystem

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type RankingSystem struct {
	RankingSystemID int64  `db:"ranking_system_id"`
	Name            string `db:"name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*RankingSystem, error) {
	d := RankingSystem{}
	err := database.Get(&d, "SELECT * FROM ranking_system where ranking_system_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]RankingSystem, error) {
	var recs []RankingSystem
	err := database.Select(&recs, "SELECT * FROM ranking_system")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *RankingSystem) error {
	_, err := database.Exec("DELETE FROM ranking_system WHERE ranking_system_id = ?", d.RankingSystemID)
	if err != nil {
		return fmt.Errorf("rankingsystem: couldn't delete rankingsystem %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *RankingSystem) error {
	res, err := database.Exec(database.BuildInsert("ranking_system", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("rankingsystem: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("rankingsystem: couldn't get last inserted ID %S", err)
	}

	d.RankingSystemID = ID

	return nil
}
