package dbuserrole

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/pkg/database"
)

type UserRole struct {
	UserRoleID int64 `db:"user_role_id"`
	UserID     int64 `db:"user_id"`
	RoleID     int64 `db:"role_id"`
}

// ReadByID reads by id column
func ReadByID(ID int64) (*UserRole, error) {
	d := UserRole{}
	err := database.Get(&d, "SELECT * FROM user_role where user_role_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ReadAll reads all records in the database
func ReadAll() ([]UserRole, error) {
	var recs []UserRole
	err := database.Select(&recs, "SELECT * FROM user_role")
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// Delete deletes a record from the database
func Delete(d *UserRole) error {
	_, err := database.Exec("DELETE FROM user_role WHERE user_role_id = ?", d.UserRoleID)
	if err != nil {
		return fmt.Errorf("userrole: couldn't delete user_role %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(d *UserRole) error {
	res, err := database.Exec(database.BuildInsert("user_role", d), database.GetArguments(*d)...)

	if err != nil {
		return fmt.Errorf("userrole: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("userrole: couldn't get last inserted ID %S", err)
	}

	d.UserRoleID = ID

	return nil
}

// ReadByUserID reads by user_id column
func ReadByUserID(ID int64) ([]UserRole, error) {
	var recs []UserRole
	err := database.Select(&recs, "SELECT * FROM user_role where user_role_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
