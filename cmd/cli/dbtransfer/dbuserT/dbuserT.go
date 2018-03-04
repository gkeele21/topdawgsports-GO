package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"topdawgsportsAPI/pkg/database/dbuser"
	"topdawgsportsAPI/pkg/database"
)

var db *sql.DB

func main() {
	// grab all users from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(127.0.0.1:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSUserID, Username, Password, DateCreated, FirstName, LastName, Email, LastLogin FROM FSUser")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userid int64
		var username, password string
		var email, firstName, lastName database.NullString
		var dateCreated, lastLogin time.Time
		if err := rows.Scan(&userid, &username, &password, &dateCreated, &firstName, &lastName, &email, &lastLogin); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("UserID : [%d], Username : [%s], Pass : [%s], Created : [%s], First : [%s], Last : [%s], Email : [%s], Last : [%s]\n", userid, username, password, dateCreated, firstName, lastName, email, lastLogin)

		user := dbuser.User{
			UserID:        userid,
			Username:      database.ToNullString(username, true),
			UserPassword:  database.ToNullString(password, true),
			CreatedDate:   dateCreated,
			FirstName:     firstName.String,
			LastName:      lastName,
			Email:         email.String,
			LastLoginDate: database.ToNullTime(lastLogin, true),
		}

		fmt.Printf("User : %#v\n", user)
		err := dbuser.Insert(&user)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
