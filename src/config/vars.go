package config

import (
	"os"
	"strconv"

	"github.com/samber/lo"
)

var Env string

var DBDialect string
var DBUrl string
var DBLogging string

var Port string

var JWTSigningKey string
var JWTExpirationDays int

var GoogleMapsAPIKey string

// should run after env.go#init as this `vars` is alphabetically after `env`
func init() {
	Env, _ = lo.Coalesce(
		os.Getenv("RAILWAY_ENVIRONMENT_NAME"),
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

	JWTSigningKey = os.Getenv("JWT_SIGNING_KEY")
	JWTExpirationDays, _ = strconv.Atoi(os.Getenv("JWT_EXPIRATION_DAYS"))

	GoogleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
}
