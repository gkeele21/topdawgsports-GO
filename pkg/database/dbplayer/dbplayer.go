package dbplayer

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type Player struct {
	PlayerID     int64               `db:"player_id"`
	SportLevelID int64               `db:"sport_level_id"`
	TeamID       database.NullInt64  `db:"team_id"`
	PositionID   database.NullInt64  `db:"position_id"`
	FirstName    database.NullString `db:"first_name"`
	LastName     database.NullString `db:"last_name"`
	Status       string              `db:"status"`
	StatsKey     database.NullString `db:"stats_key"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*Player, error) {
	d := Player{}
	err := database.Get(&d, "SELECT * FROM player where player_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]Player, error) {
	var recs []Player
	err := database.Select(&recs, "SELECT * FROM player")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *Player) error {
	_, err := database.Exec("DELETE FROM player WHERE player_id = ?", d.PlayerID)
	if err != nil {
		return fmt.Errorf("player: couldn't delete player %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *Player) error {
	res, err := database.Exec(database.BuildInsert("player", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("player: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("player: couldn't get last inserted ID %S", err)
	}

	d.PlayerID = ID

	return nil
}
