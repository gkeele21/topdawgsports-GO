package dbfantasygame

import (
	"topdawgsportsAPI/pkg/database"
	"fmt"
)

type FantasyGame struct {
	FantasyGameID  int64               `db:"fantasy_game_id"`
	GameTypeID     int64               `db:"game_type_id"`
	SportLevelID   int64               `db:"sport_level_id"`
	Name           string              `db:"name"`
	LandingPageURL database.NullString `db:"landing_page_url"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyGame, error) {
	d := FantasyGame{}
	err := database.Get(&d, "SELECT * FROM fantasy_game where fantasy_game_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyGame, error) {
	var recs []FantasyGame
	err := database.Select(&recs, "SELECT * FROM fantasy_game")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyGame) error {
	_, err := database.Exec("DELETE FROM fantasy_game WHERE fantasy_game_id = ?", d.FantasyGameID)
	if err != nil {
		return fmt.Errorf("fantasygame: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyGame) error {
	res, err := database.Exec(database.BuildInsert("fantasy_game", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasygame: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasygame: couldn't get last inserted ID %S", err)
	}

	d.FantasyGameID = ID

	return nil
}
