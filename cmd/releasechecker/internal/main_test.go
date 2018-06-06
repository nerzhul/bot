package internal

import (
	"os"
	"testing"
	"gitlab.com/nerzhul/bot/db"
)

// TestMain unit tests ramp up
func TestMain(m *testing.M) {
	gDB = &rcDB{
		config: &db.Config{
			URL:          "host=postgres dbname=unittests user=unittests password=unittests sslmode=disable",
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		},
	}

	if !gDB.init() {
		os.Exit(1)
	}

	code := m.Run()

	// Deinit code
	os.Exit(code)
}
