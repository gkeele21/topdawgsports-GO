package dbplayer

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
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
	Weight       database.NullInt64  `db:"weight"`
	Height       database.NullInt64  `db:"height"`
	Jersey       database.NullInt64  `db:"jersey"`
	Year         database.NullInt64  `db:"year"`
	HomeCity     database.NullString `db:"home_city"`
	HomeState    database.NullString `db:"home_state"`
	HomeCountry  database.NullString `db:"home_country"`
}

const STATUS_ACTIVE = "active"
const STATUS_INACTIVE = "inactive"

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

// Update will update a record in the database
func Update(s *Player) error {
	sql := database.BuildUpdate("player", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("player: couldn't update %s", err)
	}

	return nil
}

func Save(s *Player) error {
	if s.PlayerID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByStatsKeyAndSportLevel reads by stats_key and sport_level_id columns
func ReadByStatsKeyAndSportLevel(statsKey string, sportLevelId int64) (*Player, error) {
	d := Player{}
	err := database.Get(&d, "SELECT * FROM player WHERE stats_key = ? AND sport_level_id = ?", statsKey, sportLevelId)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
