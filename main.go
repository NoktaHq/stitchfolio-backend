package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imkarthi24/sf-backend/internal/app"
	"github.com/imkarthi24/sf-backend/internal/di"
)

//	@title						Stitchfolio-backend API docs
//	@version					1.0
//	@description				This is the backend for Stitchfolio.
//	@host						localhost:9000
//	@BasePath					/api/sf/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Token
//	@Security					BearerAuth
func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	var migrate string
	flag.StringVar(&migrate, migrateKey, defaultValue, usage)
	flag.Parse()

	application, err := di.InitApp(&ctx)
	checkErr(err)

	if migrate == "true" {

		application.Migrate(&ctx, checkErr)
		application.Shutdown(&ctx, checkErr)
		return
	}

	application.Start(&ctx, checkErr)

	var task *app.Task
	if application.AppConfig.Config.UseJobService {

		task, err = di.InitJobService(&ctx)
		checkErr(err)

		task.Start(&ctx, checkErr)
	}

	//Disposal
	<-done
	application.Shutdown(&ctx, checkErr)
	task.Shutdown(&ctx, checkErr)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const (
	migrateKey   = "migrate"
	defaultValue = "false"
	usage        = "Run Migration?"
)
