package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var currentInstance = new(manager)

var connectOnce sync.Once

// Connect initializes the database to be used for read/write operations
func connect() {
	connectOnce.Do(func() {
		var err error

		// Create db instance
		host := os.Getenv("TOPDAWG_DB_READ_HOST")
		username := os.Getenv("TOPDAWG_DB_USERNAME")
		password := os.Getenv("TOPDAWG_DB_PASSWORD")
		dbName := os.Getenv("TOPDAWG_DB_NAME")
		dbPort := os.Getenv("TOPDAWG_DB_PORT")

		if host == "" || username == "" || password == "" || host == "" {
			log.Fatal("invalid db config : env variables not set [TOPDAWG_DB_READ_HOST] [TOPDAWG_DB_USERNAME] [TOPDAWG_DB_PASSWORD] [TOPDAWG_DB_NAME]")
		}

		currentInstance.read, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=MST7MDT", username, password, host, dbPort, dbName))
		if err != nil {
			log.Fatal("invalid db config: ", err)
		}

		currentInstance.read.SetConnMaxLifetime(time.Second * 60)
		currentInstance.read.SetMaxIdleConns(10)
		currentInstance.read.SetMaxOpenConns(20)

		currentInstance.write = currentInstance.read

		currentInstance.WaitForConnection()
	})
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	connect()
	return currentInstance.read.QueryContext(context.Background(), query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func QueryRow(query string, args ...interface{}) *sql.Row {
	connect()
	return currentInstance.read.QueryRowContext(context.Background(), query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func Exec(query string, args ...interface{}) (sql.Result, error) {
	connect()
	return currentInstance.write.ExecContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	connect()
	return currentInstance.read.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.

// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	connect()
	return currentInstance.read.QueryRowContext(ctx, query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	connect()
	return currentInstance.write.ExecContext(ctx, query, args...)
}

// SQLX Functions

// Select scans dest array
func Select(dest interface{}, query string, args ...interface{}) error {
	connect()
	return currentInstance.read.Select(dest, query, args...)
}

// Get scans into single dest
func Get(dest interface{}, query string, args ...interface{}) error {
	connect()
	return currentInstance.read.Get(dest, query, args...)
}

// GetArguments receives a struct and puts the struct values
// into an array so that it can be passed into a query
func GetArguments(s interface{}) []interface{} {
	var args []interface{}
	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {
		args = append(args, val.Field(i).Interface())
	}

	return args
}

// GetArgumentsForUpdate receives a struct and puts the struct values
// into an array so that it can be passed into a query
func GetArgumentsForUpdate(s interface{}) []interface{} {
	var args []interface{}
	val := reflect.ValueOf(s)
	for i := 1; i < val.NumField(); i++ {
		args = append(args, val.Field(i).Interface())
	}
	args = append(args, val.Field(0).Interface())

	return args
}

// BuildInsert creates an insert query for given table name and accompanying struct
// returning a query string with placeholders ? for mysql
func BuildInsert(tableName string, obj interface{}) string {
	// build start
	query := fmt.Sprintf("INSERT INTO %s (", tableName)

	// build columns requires struct tags be set
	var columns []string
	val := reflect.ValueOf(obj)
	i := reflect.Indirect(val)
	t := i.Type()

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("db")
		if name != "" {
			columns = append(columns, name)
		}
	}

	query += strings.Join(columns, ",") + ") VALUES ("

	// build placeholders
	argsHolder := make([]string, 0)
	for range columns {
		argsHolder = append(argsHolder, "?")
	}

	query += strings.Join(argsHolder, ",")

	query += ")"

	return query
}

// BuildUpdate creates an update query for given table name and accompanying struct
// returning a query string with placeholders ? for mysql
func BuildUpdate(tableName string, obj interface{}) string {
	// build start
	query := fmt.Sprintf("UPDATE %s SET ", tableName)

	val := reflect.ValueOf(obj)
	i := reflect.Indirect(val)
	t := i.Type()

	// skip the primary key column
	for i := 1; i < t.NumField(); i++ {
		if i > 1 {
			query += ", "
		}
		name := t.Field(i).Tag.Get("db")
		if name != "" {
			query += name + " = ? "
		}
	}

	query += " WHERE " + t.Field(0).Tag.Get("db") + " = ?"

	return query
}
