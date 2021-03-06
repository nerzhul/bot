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
	log.Infof("Checking for schema migrations to perform...")
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

	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Errorf("ReleaseChecker DB Migration failed: %s", err)
		return false
	}

	log.Infof("Schema migrations done.")

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

func (db *rcDB) AddGithubRepository(group string, name string) bool {
	_, err := db.nativeDB.Exec(addGithubRepositoryQuery, group, name)
	if err != nil {
		log.Errorf("Unable to add Github repository configuration to DB: %s", err)
		return false
	}

	return true
}

func (db *rcDB) GetGithubConfiguredRepositories() ([]githubRepository, error) {
	rows, err := db.nativeDB.Query(getGithubRepositories)
	if err != nil {
		log.Errorf("Unable to execute getGithubRepositories: %s", err)
		return nil, err
	}

	defer rows.Close()

	var repoList []githubRepository

	for rows.Next() {
		gr := githubRepository{}
		if err := rows.Scan(&gr.group, &gr.name); err != nil {
			log.Errorf("Unable to read getGithubRepositories: %s", err)
			return nil, err
		}

		repoList = append(repoList, gr)
	}

	return repoList, nil
}

func (db *rcDB) RegisterRepositoryTag(group string, name string, tag string) bool {
	_, err := db.nativeDB.Exec(addGithubRepositoryTag, group, name, tag)
	if err != nil {
		log.Errorf("Unable to add Github repository tag to DB: %s", err)
		return false
	}

	return true
}

func (db *rcDB) IsGithubRepositoryTagRegistered(group string, name string, tag string) (bool, error) {
	rows, err := db.nativeDB.Query(isGithubRepositoryTagRegistered, group, name, tag)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("Unable to check if Github repository tag is registered: %s", err)
		return true, err
	}

	tagExists := false
	if rows.Next() {
		if err := rows.Scan(&tagExists); err != nil {
			log.Errorf("Unable to check if Github repository tag is registered: %s", err)
			return false, err
		}
	}

	return tagExists, nil
}

// DockerHub

func (db *rcDB) AddDockerHubImage(group string, name string) bool {
	_, err := db.nativeDB.Exec(addDockerHubImageQuery, group, name)
	if err != nil {
		log.Errorf("Unable to add DockerHub image configuration to DB: %s", err)
		return false
	}

	return true
}

func (db *rcDB) GetDockerHubConfiguredImages() ([]dockerHubImage, error) {
	rows, err := db.nativeDB.Query(getDockerHubImages)
	if err != nil {
		log.Errorf("Unable to execute getDockerHubImages: %s", err)
		return nil, err
	}

	defer rows.Close()

	var imageList []dockerHubImage

	for rows.Next() {
		gr := dockerHubImage{}
		if err := rows.Scan(&gr.group, &gr.name); err != nil {
			log.Errorf("Unable to read getDockerHubImages: %s", err)
			return nil, err
		}

		imageList = append(imageList, gr)
	}

	return imageList, nil
}

func (db *rcDB) RegisterDockerHubImageTag(group string, name string, tag string) bool {
	_, err := db.nativeDB.Exec(addDockerHubImageTag, group, name, tag)
	if err != nil {
		log.Errorf("Unable to add DockerHub image tag to DB: %s", err)
		return false
	}

	return true
}

func (db *rcDB) IsDockerHubImageTagRegistered(group string, name string, tag string) (bool, error) {
	rows, err := db.nativeDB.Query(isDockerHubImageTagRegistered, group, name, tag)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("Unable to check if DockerHub image tag is registered: %s", err)
		return true, err
	}

	tagExists := false
	if rows.Next() {
		if err := rows.Scan(&tagExists); err != nil {
			log.Errorf("Unable to check if DockerHub image tag is registered: %s", err)
			return false, err
		}
	}

	return tagExists, nil
}
