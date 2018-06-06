package internal

import (
	"os"
	"testing"
)

// TestMain unit tests ramp up
func TestMain(m *testing.M) {
	gIRCDB = &ircDB{
		config: &dbConfig{
			URL:          "host=postgres dbname=unittests user=unittests password=unittests sslmode=disable",
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		},
	}

	if !gIRCDB.init() {
		os.Exit(1)
	}

	code := m.Run()

	// Deinit code
	os.Exit(code)
}
