package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/jasonsmithj/tmp/internal/configration"
	"github.com/sirupsen/logrus"
)

type LockTable struct {
	FunctionName string
}

type DynamoDBLock interface {
	Lock()
	Unlock()
	Check() (bool, error)
}

type dynamoDBLock struct {
	db *dynamo.DB
}

func NewDynamoDBLock() *dynamoDBLock {
	return &dynamoDBLock{}
}

func (d *dynamoDBLock) Init() {
	d.db = dynamo.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
}

func (d *dynamoDBLock) Check() (bool, error) {
	table := d.db.Table(configration.Get().DynamoDBTable)
	var results = LockTable{}
	err := table.Get("FunctionName", configration.FunctionName).All(&results)
	if err != nil {
		return false, err
	}
	if len(results.FunctionName) > 0 {
		logrus.WithFields(logrus.Fields{}).Info("Multiple execution is detected and the process is terminated")
		return false, nil
	}
	return true, nil
}

func (d *dynamoDBLock) Lock() {
	table := d.db.Table(configration.Get().DynamoDBTable)

	// put item
	w := LockTable{FunctionName: configration.FunctionName}
	err := table.Put(w).Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("Failed to add lock record in dynamodb")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}
	logrus.WithFields(logrus.Fields{}).Info("Successfully added Lock records to DynamoDB")
}

func (d *dynamoDBLock) UnLock() {
	table := d.db.Table(configration.Get().DynamoDBTable)

	// delete item
	err := table.Delete("FunctionName", configration.FunctionName).Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error("Failed to add lock record in dynamodb")
		logrus.WithFields(logrus.Fields{}).Fatal(err)
	}
	logrus.WithFields(logrus.Fields{}).Info("Lock records have been successfully deleted from DynamoDB")
}
