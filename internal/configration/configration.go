package configration

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServiceAccountFile string `envconfig:"SERVICE_ACCOUNT_FILE"`
	ServiceAccountJson string `envconfig:"SERVICE_ACCOUNT_JSON"`
	GSuiteMail         string `envconfig:"GSUITE_MAIL"`
}

var globalConfig Config

// Load environment variables
func Load() {
	envconfig.MustProcess("", &globalConfig)
}

// Get environment variables
func Get() Config {
	return globalConfig
}
