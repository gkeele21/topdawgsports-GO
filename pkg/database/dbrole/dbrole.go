package dbrole

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
)

type Role struct {
	RoleID int64  `db:"role_id"`
	Name   string `db:"name"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*Role, error) {
	d := Role{}
	err := database.Get(&d, "SELECT * FROM role where role_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]Role, error) {
	var recs []Role
	err := database.Select(&recs, "SELECT * FROM role")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *Role) error {
	_, err := database.Exec("DELETE FROM role WHERE role_id = ?", d.RoleID)
	if err != nil {
		return fmt.Errorf("role: couldn't delete role %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *Role) error {
	res, err := database.Exec(database.BuildInsert("role", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("role: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("role: couldn't get last inserted ID %S", err)
	}

	d.RoleID = ID

	return nil
}
