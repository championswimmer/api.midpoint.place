package config

import (
	"os"
	"path"
	"runtime"

	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
)

func init() {
	if os.Getenv("ENV") == "test" {
		// for tests, chdir to the project root
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
		lo.Must0(os.Chdir(dir))
		if err := godotenv.Load(".env.test"); err != nil {
			applogger.Error(err)
		}
	}

	if os.Getenv("ENV") == "production" || os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
		if err := godotenv.Load(".env.production"); err != nil {
			applogger.Error(err)
		}
	}

	// Use defaults from local.env for all missing vars
	if err := godotenv.Load(".env.local"); err != nil {
		applogger.Error(err)
	}
}
