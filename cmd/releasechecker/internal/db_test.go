package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/nerzhul/bot/db"
	"os"
	"testing"
)

func TestRcDB_Init_NoConfig(t *testing.T) {
	tDB := &rcDB{}

	assert.Equal(t, false, tDB.init())
}

func TestRcDB_Init_BadURL(t *testing.T) {
	c := config{}
	c.loadDefaultConfiguration()
	c.DB.URL = "WTF"

	tDB := &rcDB{
		config: &c.DB,
	}

	assert.Equal(t, false, tDB.init())
}

func TestRcDB_Init_BadMigration(t *testing.T) {
	c := config{}
	c.loadDefaultConfiguration()

	tDB := &rcDB{
		config: &c.DB,
	}

	assert.Equal(t, false, tDB.init())
}

func TestRcDB_Init_OK(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		println("Unable to find workdir.")
		os.Exit(1)
	}

	c := config{}
	c.loadDefaultConfiguration()

	tDB := &rcDB{
		config: &db.Config{
			URL:             "host=postgres dbname=unittests user=unittests password=unittests sslmode=disable",
			MaxIdleConns:    5,
			MaxOpenConns:    10,
			MigrationSource: fmt.Sprintf("file://%s/../res/migrations", workDir),
		},
	}

	assert.Equal(t, true, tDB.init())
}

func TestRcDB_AddGithubRepository(t *testing.T) {
	assert.Equal(t, true, gDB.AddGithubRepository("nerzhul", "bot"))
}

func TestRcDB_GetGithubConfiguredRepositories(t *testing.T) {
	repositories, err := gDB.GetGithubConfiguredRepositories()
	assert.Nil(t, err)
	assert.NotNil(t, repositories)
	for _, r := range repositories {
		assert.NotEmpty(t, r.name)
		assert.NotEmpty(t, r.group)
	}
}
