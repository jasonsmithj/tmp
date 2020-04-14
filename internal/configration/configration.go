package configration

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServiceAccountJson string `envconfig:"SERVICE_ACCOUNT_JSON"`
	GSuiteMail         string `envconfig:"GSUITE_MAIL"`
}

const ServiceAccountFile = "./gsuite_admin.json"

var globalConfig Config

// Load environment variables
func Load() {
	envconfig.MustProcess("", &globalConfig)
}

// Get environment variables
func Get() Config {
	return globalConfig
}
