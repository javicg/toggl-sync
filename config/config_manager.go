package config

import (
	"github.com/spf13/viper"
	"os"
)

// Manager isolates side effects from reading and persisting config values
type Manager interface {
	Init() (ok bool, err error)
	Persist() error
}

// ViperConfigManager is an implementation of Manager that relies on "github.com/spf13/viper" for configuration management
type ViperConfigManager struct{}

// Init initializes the configuration from disk.
// If the file exists and is readable, Init returns ok=true, err=nil (after loading the configuration)
// If the file does not exist, Init returns ok=false, err=nil
// If the file is found, but cannot be read, Init returns ok=false and the error back to the client
func (mgr *ViperConfigManager) Init() (ok bool, err error) {
	viper.SetConfigName("toggl-sync")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/etc")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Persist saves the current config to disk
func (mgr *ViperConfigManager) Persist() error {
	// Creating file beforehand as viper.WriteConfig fails otherwise
	err := mgr.createConfigFile()
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}

func (mgr *ViperConfigManager) createConfigFile() error {
	f, err := os.OpenFile("/usr/local/etc/toggl-sync.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	return f.Close()
}
