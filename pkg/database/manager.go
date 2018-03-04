package database

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

// manager is used to manage the read and write database connections
type manager struct {
	read  *sqlx.DB
	write *sqlx.DB
}

// Ping checks both read and write db for connection
func (manager *manager) Ping() error {
	err := manager.read.Ping()
	if err != nil {
		return err
	}

	err = manager.write.Ping()
	if err != nil {
		return err
	}

	return nil
}

// WaitForConnection waits 30 seconds for connection before panic by pinging every 2 seconds
func (manager *manager) WaitForConnection() {
	tryCount := 0
	for {
		// Attempt connection
		err := manager.Ping()
		if err == nil {
			break
		}

		time.Sleep(2 * time.Second)
		tryCount++

		if tryCount == 15 {
			log.Println("db connection not available within 30s")
			os.Exit(1)
		}
	}
}
