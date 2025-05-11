package config

import (
	"os"

	"github.com/samber/lo"
)

var Env string

var DBDialect string
var DBUrl string
var DBLogging string

var Port string

// should run after env.go#init as this `vars` is alphabetically after `env`
func init() {
	Env, _ = lo.Coalesce(
		os.Getenv("RAILWAY_ENVIRONMENT"),
		os.Getenv("ENV"),
		"local",
	)

	DBDialect = os.Getenv("DB_DIALECT")
	DBUrl, _ = lo.Coalesce(
		os.Getenv("DATABASE_PRIVATE_URL"),
		os.Getenv("DATABASE_URL"),
	)

	DBLogging = os.Getenv("DB_LOGGING")

	Port = os.Getenv("PORT")
}
