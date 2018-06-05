package internal

import (
	"database/sql"
	_ "github.com/lib/pq" // pq requires blank import
	dblib "gitlab.com/nerzhul/bot/db"
)

type rcDB struct {
	nativeDB *sql.DB
	config   *dblib.Config
}

func (db *rcDB) init(config *dblib.Config) bool {
	log.Infof("Connecting to IRC DB at %s", config.URL)
	nativeDB, err := sql.Open("postgres", config.URL)
	if err != nil {
		log.Errorf("Failed to connect to IRC DB: %s", err)
		return false
	}

	db.nativeDB = nativeDB
	if !db.ValidationQuery() {
		db.nativeDB = nil
		return false
	}

	db.nativeDB.SetMaxIdleConns(config.MaxIdleConns)
	db.nativeDB.SetMaxOpenConns(config.MaxOpenConns)

	log.Infof("Connected to IRC DB.")
	return true
}

func (db *rcDB) ValidationQuery() bool {
	rows, err := db.nativeDB.Query(dblib.ValidationQuery)
	if err != nil {
		log.Errorf("Failed to run IRC validation query: %s", err)
		return false
	}
	rows.Close()
	return true
}
