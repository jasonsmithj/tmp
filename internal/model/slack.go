package model

import (
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/jasonsmithj/tmp/internal/configration"
	"github.com/sirupsen/logrus"
)

type SlackNotification interface {
	Send(user string, password string)
}

type slackNotification struct{}

func NewSlackNotification() SlackNotification {
	return &slackNotification{}
}

func (s *slackNotification) Send(user string, password string) {
	payload := slack.Payload{
		Text:      fmt.Sprintf("Change Password\nUser Name: %s\nPassword: %s", user, password),
		Username:  "robot",
		Channel:   "#slack_test",
		IconEmoji: ":monkey_face:",
	}
	err := slack.Send(configration.Get().WebHookUrl, "", payload)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("The slack notification failed")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}
	logrus.WithFields(logrus.Fields{
		"user":     user,
		"password": password,
	}).Info("Successfully changed SLACK notifications and passwords")
}
