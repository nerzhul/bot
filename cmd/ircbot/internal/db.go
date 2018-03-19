package internal

import (
	"database/sql"
	_ "github.com/lib/pq" // pq requires blank import
)

type dbConfig struct {
	URL          string `yaml:"url"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
}

type ircDB struct {
	nativeDB *sql.DB
	config   *dbConfig
}

func (db *ircDB) init(config *dbConfig) bool {
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

func (db *ircDB) ValidationQuery() bool {
	rows, err := db.nativeDB.Query(ValidationQuery)
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
		return nil, nil
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
