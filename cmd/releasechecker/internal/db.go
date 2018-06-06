package internal

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"   // golang-migrate requires blank import
	_ "github.com/golang-migrate/migrate/source/github" // golang-migrate requires blank import
	_ "github.com/lib/pq"                               // pq requires blank import
	dblib "gitlab.com/nerzhul/bot/db"
)

type rcDB struct {
	nativeDB *sql.DB
	config   *dblib.Config
}

func (db *rcDB) init() bool {
	if db.config == nil {
		log.Errorf("DB config is nil !!!")
		return false
	}

	log.Infof("Connecting to ReleaseChecker DB at %s", db.config.URL)
	nativeDB, err := sql.Open("postgres", db.config.URL)
	if err != nil {
		log.Errorf("Failed to connect to ReleaseChecker DB: %s", err)
		return false
	}

	db.nativeDB = nativeDB
	if !db.ValidationQuery() {
		db.nativeDB = nil
		return false
	}

	db.nativeDB.SetMaxIdleConns(db.config.MaxIdleConns)
	db.nativeDB.SetMaxOpenConns(db.config.MaxOpenConns)

	if !db.runMigrations() {
		db.nativeDB = nil
		return false
	}

	log.Infof("Connected to ReleaseChecker DB.")
	return true
}

func (db *rcDB) runMigrations() bool {
	driver, err := postgres.WithInstance(db.nativeDB, &postgres.Config{})
	if err != nil {
		log.Errorf("Unable to create migration instance on ReleaseChecker DB: %s", err)
		return false
	}

	m, err := migrate.NewWithDatabaseInstance(
		db.config.MigrationSource,
		"postgres", driver)

	if err != nil {
		log.Errorf("Unable to run migrations on ReleaseChecker DB: %s", err)
		return false
	}

	m.Steps(2)
	return true
}

func (db *rcDB) ValidationQuery() bool {
	rows, err := db.nativeDB.Query(dblib.ValidationQuery)
	if err != nil {
		log.Errorf("Failed to run ReleaseChecker validation query: %s", err)
		return false
	}
	rows.Close()
	return true
}
