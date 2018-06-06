package internal

import (
	"fmt"
	"gitlab.com/nerzhul/bot/db"
	"os"
	"testing"
)

// TestMain unit tests ramp up
func TestMain(m *testing.M) {
	workDir, err := os.Getwd()
	if err != nil {
		println("Unable to find workdir.")
		os.Exit(1)
	}
	gDB = &rcDB{
		config: &db.Config{
			URL:             "host=postgres dbname=unittests user=unittests password=unittests sslmode=disable",
			MaxIdleConns:    5,
			MaxOpenConns:    10,
			MigrationSource: fmt.Sprintf("file://%s/../res/migrations", workDir),
		},
	}

	if !gDB.init() {
		os.Exit(1)
	}

	// Reinit some data
	gDB.nativeDB.Exec("TRUNCATE TABLE github_repository_tags")

	code := m.Run()

	// Deinit code
	os.Exit(code)
}
