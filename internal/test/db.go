package test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // sql driver
	"github.com/redhajuanda/gorengan/config"
	migrate "github.com/rubenv/sql-migrate"
)

var db *sql.DB
var err error

// DB returns a new test DB
func newDB(t *testing.T) *sql.DB {
	cfg := config.LoadTest()
	connString := fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local&multiStatements=true", cfg.Database.Username, cfg.Database.Password, "")
	// connect DB
	db, err = sql.Open("mysql", connString)
	if err != nil {
		t.Errorf("Error opening Test DB: %v", err)
		t.FailNow()
	}

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %v; CREATE DATABASE IF NOT EXISTS %v; USE %v", cfg.Database.DBName, cfg.Database.DBName, cfg.Database.DBName)

	_, err = db.Exec(query)
	if err != nil {
		t.Errorf("Error preparing database: %v", err)
		t.FailNow()
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "../../migrations",
	}
	migrated, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		t.Errorf("Error migrating database: %v", err)
		t.FailNow()
	}
	fmt.Printf("%v migrations applied", migrated)

	return db
}

// GetTestDB returns a DB connection
func GetTestDB(t *testing.T) *sql.DB {
	if db == nil {
		fmt.Println("======== creating test db ========")
		return newDB(t)
	}
	return db
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *sql.DB, tables ...string) {
	fmt.Println("======== truncate table ===========")
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %v", table))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
