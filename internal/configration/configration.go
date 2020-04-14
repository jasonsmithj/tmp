package configration

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServiceAccountJson string `envconfig:"SERVICE_ACCOUNT_JSON"`
	GSuiteMail         string `envconfig:"GSUITE_MAIL"`
	WebHookUrl         string `envconfig:"WEB_HOOK_URL"`
	DynamoDBTable      string `envconfig:"DYNAMO_DB_TABLE"`
}

const (
	ServiceAccountFile = "/tmp/gsuite_admin.json"
	FunctionName       = "gsuitePasswordChange"
)

var globalConfig Config

// Load environment variables
func Load() {
	envconfig.MustProcess("", &globalConfig)
}

// Get environment variables
func Get() Config {
	return globalConfig
}
