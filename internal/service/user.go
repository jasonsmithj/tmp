package service

import (
	"math/rand"
	"time"

	"github.com/jasonsmithj/tmp/internal/model"
	"github.com/sirupsen/logrus"
	admin "google.golang.org/api/admin/directory/v1"
)

type GSuiteUser interface {
	Get(service *admin.Service) *admin.Users
	Update(service *admin.Service, users *admin.Users)
}

type gsuiteUser struct{}

func NewGSuiteUser() GSuiteUser {
	return &gsuiteUser{}
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	rs6Letters       = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func (g *gsuiteUser) Get(service *admin.Service) *admin.Users {
	r, err := service.Users.List().Customer("my_customer").MaxResults(100).Query("email:takayuki*").Do()
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("Failed to retrieve user information")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}
	if len(r.Users) == 0 {
		logrus.WithFields(logrus.Fields{}).Fatal("No users found")
	}
	return r
}

func (g *gsuiteUser) Update(service *admin.Service, users *admin.Users) {
	newSlackNotification := model.NewSlackNotification()
	for _, u := range users.Users {
		user := u.PrimaryEmail
		password := GeneratePassword(16)
		u.Password = password
		_, err := service.Users.Update(u.PrimaryEmail, u).Do()
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error("Failed to change password")
			logrus.WithFields(logrus.Fields{}).Fatal(err)
		}
		logrus.WithFields(logrus.Fields{
			"user":     u.PrimaryEmail,
			"password": password,
		}).Info("Password changed successfully!")
		newSlackNotification.Send(user, password)
	}
}

func GeneratePassword(n int) string {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return string(b)
}
