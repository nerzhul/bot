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
	assert.Equal(t, true, gDB.AddGithubRepository("minetest", "minetest"))
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

func TestRcDB_RegisterRepositoryTag(t *testing.T) {
	assert.Equal(t, true, gDB.RegisterRepositoryTag("minetest", "minetest", "0.4.15"))
}

func TestRcDB_IsGithubRepositoryTagRegistered(t *testing.T) {
	exists, err := gDB.IsGithubRepositoryTagRegistered("minetest", "minetest", "0.4.15")
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	exists, err = gDB.IsGithubRepositoryTagRegistered("testnonexistent", "testnonexistent", "0.0.0")
	assert.Nil(t, err)
	assert.Equal(t, false, exists)
}

func TestRcDB_AddDockerHubImage(t *testing.T) {
	assert.Equal(t, true, gDB.AddDockerHubImage("arm64v8", "httpd"))
}

func TestRcDB_GetDockerHubConfiguredImages(t *testing.T) {
	images, err := gDB.GetDockerHubConfiguredImages()
	assert.Nil(t, err)
	assert.NotNil(t, images)
	for _, r := range images {
		assert.NotEmpty(t, r.name)
		assert.NotEmpty(t, r.group)
	}
}

func TestRcDB_RegisterDockerHubImageTag(t *testing.T) {
	assert.Equal(t, true, gDB.RegisterDockerHubImageTag("arm64v8", "httpd", "2.4.1"))
}

func TestRcDB_IsDockerHubImageTagRegistered(t *testing.T) {
	exists, err := gDB.IsDockerHubImageTagRegistered("arm64v8", "httpd", "2.4.1")
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	exists, err = gDB.IsDockerHubImageTagRegistered("testnonexistent", "testnonexistent", "0.0.0")
	assert.Nil(t, err)
	assert.Equal(t, false, exists)
}
