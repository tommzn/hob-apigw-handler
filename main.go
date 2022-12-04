package main

import (
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	timetracker "github.com/tommzn/hob-timetracker"
)

func main() {

	handler, err := bootstrap()
	if err != nil {
		panic(err)
	}
	lambda.Start(handler.Process)
}

// bootstrap loads config and creates a new scheduled collector with a exchangerate datasource.
func bootstrap() (Handler, error) {

	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	secretsManager := newSecretsManager()
	logger := newLogger(conf, secretsManager)
	timeTracker, err := newTimeTracker(conf)
	if err != nil {
		return nil, err
	}
	return newRequestHandler(timeTracker, logger), nil
}

// loadConfig from config file.
func loadConfig() (config.Config, error) {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		return nil, err
	}
	return configSource.Load()
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	return log.WithNameSpace(logger, "hob-apigw-handler")
}

// NewTimeTracker creates a new time tracker to persist records in a S3 bucket.
func newTimeTracker(conf config.Config) (timetracker.TimeTracker, error) {

	awsRegion := conf.Get("aws.s3.region", config.AsStringPtr(os.Getenv("AWS_REGION")))
	bucket := conf.Get("aws.s3.bucket", nil)
	if bucket == nil {
		return nil, errors.New("Np S3 bucket specified!")
	}
	basePath := conf.Get("aws.s3.basepath", nil)
	return timetracker.NewS3Repository(awsRegion, bucket, basePath), nil
}
