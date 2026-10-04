package config

import (
	"github.com/spf13/viper"
)

// Load reads yaml/json config files from the given paths and returns a *viper.Viper.
func Load(configFiles ...string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	for _, f := range configFiles {
		v.AddConfigPath(".")
		v.SetConfigName(f) // viper will look for f.yaml or f.yml
	}
	if err := v.ReadInConfig(); err != nil {
		// If the file does not exist we still return a viper that can be overridden
		// by env vars or flags – this lets the application start with defaults.
	}
	v.AutomaticEnv() // allow overriding via env vars (e.g., DB_PASSWORD)
	return v
}
