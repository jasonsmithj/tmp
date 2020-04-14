package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jasonsmithj/tmp/internal/configration"
	"github.com/jasonsmithj/tmp/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {
	// Load and set the environment variables.
	configration.Load()

	// Sets the format of the log output.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Using DynamoDB to prevent multiple execution of Lambda
	dynamodbLock := service.NewDynamoDBLock()
	dynamodbLock.Init()
	r, err := dynamodbLock.Check()
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("Failed to get a lock record from DynamoDB")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}
	if r == false {
		return
	}
	dynamodbLock.Lock()

	directoryService := service.NewDirectoryService()
	// Auth a json file from the environment variables
	directoryService.CreateKey()
	// Authentication of GSuite is performed
	adminService, err := directoryService.Auth(configration.Get().GSuiteMail)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("google authentication failed")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}

	// Get GSuite Users
	gsuiteUser := service.NewGSuiteUser()
	u := gsuiteUser.Get(adminService)
	gsuiteUser.Update(adminService, u)

	dynamodbLock.UnLock()
}
