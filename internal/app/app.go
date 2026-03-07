package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	config_cache "github.com/imkarthi24/sf-backend/internal/cache"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository"

	// "github.com/loop-kar/pixie/db/migrator"

	"github.com/loop-kar/pixie/log"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/loop-kar/pixie/db/migrator"

	"gorm.io/gorm"
)

type App struct {
	Server                 *gin.Engine
	AppConfig              config.AppConfig
	StitchDB               *gorm.DB
	NewRelic               *newrelic.Application
	MasterConfigRepository repository.MasterConfigRepository
}

func (a *App) Start(ctx *context.Context, checkErr func(err error)) {
	go func(ctx *context.Context) {

		//App startup essentials
		config := log.PlogConfig{
			Level:            a.AppConfig.Logger.LogLevel,
			Environment:      a.AppConfig.Server.Environment,
			OutputPaths:      []string{"stdout"},
			EnableCaller:     true,
			EnableStacktrace: true,
			Provider:         a.AppConfig.Logger.Provider,
		}
		// Initialize with error handling
		if err := log.InitLoggerWithConfig(a.NewRelic, config); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		defer log.Sync()

		log.Info(ctx, "App Service Started...")
		config_cache.InitMasterConfig(a.MasterConfigRepository)

		host := a.AppConfig.Server.Host
		if strings.Contains(host, "render") {
			host = ""
		}
		address := fmt.Sprintf("%s:%d", host, a.AppConfig.Server.Port)
		err := a.Server.Run(address)
		if err != nil {
			checkErr(err)
		}
	}(ctx)
}

func (a *App) Shutdown(ctx *context.Context, checkErr func(err error)) {
	dbConn, err := a.StitchDB.DB()
	checkErr(err)

	err = dbConn.Close()
	checkErr(err)

	log.Info(ctx, "App Service Shutdown...")
	log.Sync()
}

func (a *App) Migrate(ctx *context.Context, checkErr func(err error)) {
	log.MustInitLogger(a.NewRelic)
	migrator := migrator.NewMigrator(a.StitchDB)

	entityList := []interface{}{
		// &entities.Channel{},
		// &entities.Customer{},
		// &entities.DressType{},
		// &entities.EmailNotification{},
		// &entities.EnquiryHistory{},
		// &entities.Enquiry{},
		// &entities.Expense{},
		// &entities.MasterConfig{},
		// &entities.Measurement{},
		// &entities.MeasurementHistory{},
		// &entities.Notification{},
		// &entities.OrderHistory{},
		// &entities.Order{},
		// &entities.OrderItem{},
		// &entities.Person{},
		// &entities.Task{},
		// &entities.UserChannelDetail{},
		// &entities.UserConfig{},
		// &entities.User{},
		// &entities.WhatsappNotification{},
		//&entities.Task{},
		// &entities.Inventory{},
		// &entities.InventoryLog{},
		// &entities.Product{},
		// &entities.Category{},
		&entities.FileStoreMetadata{},
	}

	//************************//
	/*
	 0. Use Migrate only for initial migration
	 1. Make changes to entity files
	 2. Provide proper migration name in GenerateAlterMigration with incremented number and snake case eg:" 002_person_entity_update"
	 4. Only uncomment entity which is changed in entityList to generate minimal alter migration file
	 5. Comment back the GenerateAlterMigration line after generating migration file to avoid overwriting
	 6. Use GenerateAlterMigration to generate alter migration file
	*/
	//************************//

	//migrator.Migrate(entityList, checkErr)

	migrator.GenerateAlterMigration(entityList, "009_file_store_entity_update")
}
