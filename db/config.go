package db

type Config struct {
	URL             string `yaml:"url"`
	MaxIdleConns    int    `yaml:"max-idle-conns"`
	MaxOpenConns    int    `yaml:"max-open-conns"`
	MigrationSource string `yaml:"db-migration-source"`
}
