package db

import (
	"sync"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var appDB *gorm.DB
var createAppDbOnce sync.Once

func getGormConfig() (dbConfig *gorm.Config) {
	dbConfig = &gorm.Config{
		TranslateError: true,
	}
	var dbLogLevel logger.LogLevel = lo.Switch[string, logger.LogLevel](config.DBLogging).
		Case("info", logger.Info).
		Case("warn", logger.Warn).
		Case("error", logger.Error).
		Default(logger.Error)

	dbConfig.Logger = logger.Default.LogMode(dbLogLevel)
	return
}

type DatabaseProvider func(dbUrl string, config *gorm.Config) *gorm.DB

var dbProviders map[string]DatabaseProvider = map[string]DatabaseProvider{}

func InjectDBProvider(name string, provider DatabaseProvider) {
	dbProviders[name] = provider
}

func init() {
	InjectDBProvider("sqlite", ProvideSqliteDB)
	InjectDBProvider("postgres", ProvidePostgresDB)
}

func GetAppDB() *gorm.DB {

	createAppDbOnce.Do(func() {
		applogger.Warn("App: Initialising database")
		switch config.DBDialect {
		case "sqlite":
			appDB = dbProviders["sqlite"](config.DBUrl, getGormConfig())
		case "postgres":
			appDB = dbProviders["postgres"](config.DBUrl, getGormConfig())
		default:
			panic("Database config incorrect")
		}

		lo.Must0(appDB.AutoMigrate(&models.User{}))
		lo.Must0(appDB.AutoMigrate(&models.Group{}))
		lo.Must0(appDB.AutoMigrate(&models.GroupUser{}))
		lo.Must0(appDB.AutoMigrate(&models.GroupPlace{}))
		lo.Must0(appDB.AutoMigrate(&models.WaitlistSignup{}))

	})

	return appDB
}
