package dbuser

import (
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"time"
)

type User struct {
	UserID        int64               `db:"user_id"`
	Email         string              `db:"email"`
	Username      database.NullString `db:"username"`
	UserPassword  database.NullString `db:"user_password"`
	FirstName     string              `db:"first_name"`
	LastName      database.NullString `db:"last_name"`
	Cell          database.NullString `db:"cell"`
	CreatedDate   time.Time           `db:"created_date"`
	LastLoginDate database.NullTime   `db:"last_login_date"`
}

// ReadByID reads user by id column
func ReadByID(ID int64) (*User, error) {
	u := User{}
	err := database.Get(&u, "SELECT * FROM user where user_id = ?", ID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// ReadAll reads all users in the database
func ReadAll() ([]User, error) {
	var users []User
	err := database.Select(&users, "SELECT * FROM user")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Delete deletes a user from the database
func Delete(u *User) error {
	_, err := database.Exec("DELETE FROM user WHERE user_id = ?", u.UserID)
	if err != nil {
		return fmt.Errorf("user: couldn't delete user %s", err)
	}

	return nil
}

// Insert will create a new record in the database
func Insert(u *User) error {
	res, err := database.Exec(database.BuildInsert("user", u), database.GetArguments(*u)...)

	if err != nil {
		return fmt.Errorf("user: couldn't insert new %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("user: couldn't get last inserted ID %S", err)
	}

	u.UserID = ID

	return nil
}

// Update will update a record in the database
func Update(s *User) error {
	sql := database.BuildUpdate("user", s)
	_, err := database.Exec(sql, database.GetArgumentsForUpdate(*s)...)

	if err != nil {
		return fmt.Errorf("user: couldn't update %s", err)
	}

	return nil
}

func Save(s *User) error {
	if s.UserID > 0 {
		return Update(s)
	} else {
		return Insert(s)
	}
}

// ReadByUsername reads user by username column
func ReadByUsername(username string) (*User, error) {
	u := User{}
	err := database.Get(&u, "SELECT * FROM user where username = ?", username)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
