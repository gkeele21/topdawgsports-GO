package dbfantasyleague

import (
	"fmt"
	"time"
	"topdawgsportsAPI/pkg/database"
)

type FantasyLeague struct {
	FantasyLeagueID int64               `db:"fantasy_league_id"`
	SeasonID        int64               `db:"season_id"`
	FantasyGameID   int64               `db:"fantasy_game_id"`
	Name            string              `db:"name"`
	Description     database.NullString `db:"description"`
	LeaguePassword  database.NullString `db:"league_password"`
	Visibility      string              `db:"visibility"`
	CreatedDate     time.Time           `db:"created_date"`
	CreatedByUserID database.NullInt64  `db:"created_by_user_id"`
	Status          string              `db:"status"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*FantasyLeague, error) {
	d := FantasyLeague{}
	err := database.Get(&d, "SELECT * FROM fantasy_league where fantasy_league_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]FantasyLeague, error) {
	var recs []FantasyLeague
	err := database.Select(&recs, "SELECT * FROM fantasy_league")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *FantasyLeague) error {
	_, err := database.Exec("DELETE FROM fantasy_league WHERE fantasy_league_id = ?", d.FantasyLeagueID)
	if err != nil {
		return fmt.Errorf("fantasyleague: couldn't delete record %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *FantasyLeague) error {
	res, err := database.Exec(database.BuildInsert("fantasy_league", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("fantasyleague: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fantasyleague: couldn't get last inserted ID %S", err)
	}

	d.FantasyLeagueID = ID

	return nil
}

// Update will update a record in the database
func Update(s *FantasyLeague) error {
	sql := database.BuildUpdate("fantasy_league", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("season: couldn't update %s", err)
	}

	return nil
}


// ReadAllBySeasonID_FantasyGameID reads all fantasy_leagues in the database for the given seasonID and gameID
func ReadAllBySeasonIDFantasyGameID(seasonID, gameID int64, orderBy string) ([]FantasyLeague, error) {
	var recs []FantasyLeague
	if orderBy == "" {
		orderBy = "fantasy_league_id asc"
	}
	err := database.Select(&recs, "SELECT * FROM fantasy_league WHERE season_id = ? AND fantasy_game_id = ? ORDER BY "+orderBy, seasonID, gameID)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
