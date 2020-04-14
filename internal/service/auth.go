package service

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jasonsmithj/tmp/internal/configration"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type DirectoryService interface {
	Auth(userEmail string) (*admin.Service, error)
	CreateKey()
}

type directoryService struct{}

func NewDirectoryService() DirectoryService {
	return &directoryService{}
}

func (d *directoryService) Auth(userEmail string) (*admin.Service, error) {
	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile(configration.ServiceAccountFile)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryUserScope)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}
	config.Subject = userEmail

	ts := config.TokenSource(ctx)

	srv, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("NewService: %v", err)
	}
	return srv, nil
}

func (d *directoryService) CreateKey() {
	_, err := os.Stat(configration.ServiceAccountFile)
	if os.IsNotExist(err) {
		keyPayload, err := base64.StdEncoding.DecodeString(configration.Get().ServiceAccountJson)

		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error("Failed to decode the service account key from base64")
			logrus.WithFields(logrus.Fields{}).Fatal(err)
		}

		file, err := os.Create(configration.ServiceAccountFile)
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error("Failed to create key file for service account")
			logrus.WithFields(logrus.Fields{}).Fatal(err)
		}
		defer file.Close()

		_, err = file.Write(([]byte)(keyPayload))
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error("Failed to write key file for service account")
			logrus.WithFields(logrus.Fields{}).Fatal(err)
		}
		logrus.WithFields(logrus.Fields{}).Info("Successfully created a key file for the service account")
	} else {
		logrus.WithFields(logrus.Fields{}).Info("The key file for the service account exists.")
	}

}
