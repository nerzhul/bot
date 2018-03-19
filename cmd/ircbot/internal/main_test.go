package internal

import (
	"os"
	"testing"
)

// TestMain unit tests ramp up
func TestMain(m *testing.M) {
	gIRCDB = &ircDB{}

	if !gIRCDB.init(&dbConfig{
		URL:          "host=postgres dbname=unittests user=unittests password=unittests sslmode=disable",
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}) {
		os.Exit(1)
	}

	code := m.Run()

	// Deinit code
	os.Exit(code)
}
