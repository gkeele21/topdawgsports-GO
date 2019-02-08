package dbteam

import (
	"fmt"
	"topdawgsportsAPI/pkg/database"
)

type Team struct {
	TeamID       int64               `db:"team_id"`
	Name         string              `db:"name"`
	Abbreviation database.NullString `db:"abbreviation"`
	Mascot       database.NullString `db:"mascot"`
}

// ReadByID reads user by id column
func ReadByID(ID int64) (*Team, error) {
	t := Team{}
	err := database.Get(&t, "SELECT * FROM team where team_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// ReadAll reads all teams in the database
func ReadAll() ([]Team, error) {
	var teams []Team
	err := database.Select(&teams, "SELECT * FROM team")
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// Delete deletes a team from the database
func Delete(t *Team) error {
	_, err := database.Exec("DELETE FROM team WHERE team_id = ?", t.TeamID)
	if err != nil {
		return fmt.Errorf("team: couldn't delete team %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(t *Team) error {
	res, err := database.Exec(database.BuildInsert("team", t), database.GetArguments(*t)...)

	if err != nil {
		return fmt.Errorf("team: couldn't insert new %s %#v", err, t)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("team: couldn't get last inserted ID %S", err)
	}

	t.TeamID = ID

	return nil
}
