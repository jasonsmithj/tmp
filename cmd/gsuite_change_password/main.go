package main

import (
	"github.com/jasonsmithj/tmp/internal/configration"
	"github.com/jasonsmithj/tmp/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load and set the environment variables.
	configration.Load()

	// Sets the format of the log output.
	logrus.SetFormatter(&logrus.JSONFormatter{})


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
	r := gsuiteUser.Get(adminService)
	gsuiteUser.Update(adminService, r)
}
