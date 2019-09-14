package dbfantasygametype

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
)

type FantasyGameType struct {
	FantasyGameTypeID int64               `db:"fantasy_game_type_id"`
	Name              string              `db:"name"`
	DisplayName       database.NullString `db:"display_name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyGameType, error) {
	d := FantasyGameType{}
	err := database.Get(&d, "SELECT * FROM fantasy_game_type_id where fantasy_game_type_id_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyGameType, error) {
	var recs []FantasyGameType
	err := database.Select(&recs, "SELECT * FROM fantasy_game_type_id")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyGameType) error {
	_, err := database.Exec("DELETE FROM fantasy_game_type_id WHERE fantasy_game_type_id_id = ?", d.FantasyGameTypeID)
	if err != nil {
		return fmt.Errorf("fantasygametype: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyGameType) error {
	res, err := database.Exec(database.BuildInsert("fantasy_game_type", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasygametype: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasygametype: couldn't get last inserted ID %S", err)
	}

	d.FantasyGameTypeID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyGameType) error {
	sql := database.BuildUpdate("fantasy_game_type", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("fantasy_game_type: couldn't update %s", err)
	}

	return nil
}

func Save(s *FantasyGameType) error {
	if s.FantasyGameTypeID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}
