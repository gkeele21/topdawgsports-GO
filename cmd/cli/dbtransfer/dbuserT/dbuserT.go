package main

import (
	"database/sql"
	"fmt"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database"
	"github.com/gkeele21/topdawgsportsAPI/internal/app/database/dbuser"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB

func main() {
	// grab all users from the existing database
	//db, err := sql.Open("mysql", "webuser:lakers55@tcp(topdawg.circlepix.com:3306)/topdawg?parseTime=true")
	db, err := sql.Open("mysql", "sumo:password@tcp(localhost:3307)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT FSUserID, Username, Password, DateCreated, FirstName, LastName, Email, LastLogin FROM fsuser WHERE FSUserID > 2")
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

		if userid < 166 {
			continue;
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
