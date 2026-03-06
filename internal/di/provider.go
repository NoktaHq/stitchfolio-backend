package di

import (
	"fmt"
	"os"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/storage"
	"github.com/loop-kar/pixie/storage/client"
	"github.com/loop-kar/pixie/storage/s3"
	pkgservice "github.com/loop-kar/pixie/service"
	pkgemail "github.com/loop-kar/pixie/service/email"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// ProvideServiceContainer builds the shared "pkg/service" dependency container.
// This is intended to be injected into internal services via Wire.
func ProvideServiceContainer(appConfig config.AppConfig) *pkgservice.Service {
	emailSvc := pkgemail.NewEmailService(pkgemail.SMTPConfig{
		UserName:   appConfig.SMTP.UserName,
		Password:   appConfig.SMTP.Password,
		Host:       appConfig.SMTP.Host,
		Port:       appConfig.SMTP.Port,
		Override:   appConfig.SMTP.Override,
		OverrideTo: appConfig.SMTP.OverrideTo,
	})

	return pkgservice.NewService(
		pkgservice.WithEmailService(emailSvc),
	)
}

// ProvideDatabaseConnectionParams maps the internal config to the database connection params
func ProvideDatabaseConnectionParams(appConfig config.AppConfig) db.DatabaseConnectionParams {
	sslMode := "prefer"

	return db.DatabaseConnectionParams{
		Host:        appConfig.Database.Host,
		Port:        appConfig.Database.Port,
		Username:    appConfig.Database.Username,
		DBName:      appConfig.Database.DBName,
		Password:    appConfig.Database.Password,
		SSLMode:     sslMode,
		Schema:      appConfig.Database.Schema,
		Environment: appConfig.Server.Environment,
	}
}

// ProvideNewRelic initializes the New Relic application with the provided configuration
func ProvideNewRelic(appConfig config.AppConfig) *newrelic.Application {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appConfig.Server.AppName),
		newrelic.ConfigLicense(appConfig.Logger.License),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if nil != err {
		fmt.Printf("New Relic initialization failed: %v\n", err)
		os.Exit(1)
	}

	return app
}

// ProvideCloudStorageProvider returns a CloudStorageProvider from the app's S3 config (for file-store and entity-document flows).
func ProvideCloudStorageProvider(appConfig config.AppConfig) (storage.CloudStorageProvider, error) {
	cfg := s3.S3Config{
		Region:          appConfig.S3Config.Region,
		Bucket:          appConfig.S3Config.Bucket,
		AccessKeyID:     appConfig.S3Config.AccessKeyID,
		SecretAccessKey: appConfig.S3Config.SecretAccessKey,
		Endpoint:        appConfig.S3Config.Endpoint,
		UsePathStyle:    appConfig.S3Config.UsePathStyle,
		ForceHTTPS:      appConfig.S3Config.ForceHTTPS,
	}
	return client.ProvideCloudStorageClient(cfg)
}
