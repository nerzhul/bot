package internal

import (
	"database/sql"
	_ "github.com/lib/pq" // pq requires blank import
	dblib "gitlab.com/nerzhul/bot/db"
)

type ircDB struct {
	nativeDB *sql.DB
	config   *dblib.Config
}

func (db *ircDB) init() bool {
	if db.config == nil {
		log.Fatalf("DB config is nil, aborting !")
	}

	log.Infof("Connecting to IRC DB at %s", db.config.URL)
	nativeDB, err := sql.Open("postgres", db.config.URL)
	if err != nil {
		log.Errorf("Failed to connect to IRC DB: %s", err)
		return false
	}

	db.nativeDB = nativeDB
	if !db.ValidationQuery() {
		db.nativeDB = nil
		return false
	}

	db.nativeDB.SetMaxIdleConns(db.config.MaxIdleConns)
	db.nativeDB.SetMaxOpenConns(db.config.MaxOpenConns)

	log.Infof("Connected to IRC DB.")
	return true
}

func (db *ircDB) ValidationQuery() bool {
	rows, err := db.nativeDB.Query(dblib.ValidationQuery)
	if err != nil {
		log.Errorf("Failed to run IRC validation query: %s", err)
		return false
	}
	rows.Close()
	return true
}

func (db *ircDB) loadIRCChannelConfigs() ([]ircChannelConfig, error) {
	rows, err := db.nativeDB.Query(LoadChannelConfigQuery)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("Error while loading IRC channel configurations: %s", err)
		return nil, err
	}

	var configList []ircChannelConfig

	for rows.Next() {
		config := ircChannelConfig{}
		if err := rows.Scan(&config.Name, &config.Password, &config.AnswerCommands, &config.Hello); err != nil {
			log.Errorf("Error while scanning IRC channel configurations: %s", err)
			return nil, err
		}

		configList = append(configList, config)
	}

	return configList, nil
}

// SaveIRCChannelConfig save irc channel configuration for name channel with optional password
func (db *ircDB) SaveIRCChannelConfig(name string, password string) error {
	log.Debugf("Saving IRCChannel configuration for channel '%s'", name)
	rows, err := db.nativeDB.Query(RegisterChannelConfigQuery, name, password)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("Error while saving IRC channel configurations: %s", err)
		return err
	}

	return nil
}
