package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"topdawgsportsAPI/pkg/database/dbdivision"
)

func main() {
	// grab all teams from the existing database
	db, err := sql.Open("mysql", "webuser:lakers55@tcp(topdawg.circlepix.com:3306)/topdawg?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT DivisionID, DisplayName FROM Division")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var divisionid int64
		var name string
		if err := rows.Scan(&divisionid, &name); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("DivisionID : [%d], Name : [%s]\n", divisionid, name)

		div := dbdivision.Division{
			DivisionID: divisionid,
			Name:       name,
		}

		fmt.Printf("Record : %#v\n", div)
		err := dbdivision.Insert(&div)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
