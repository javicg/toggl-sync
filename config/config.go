package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// FileUsed returns the absolute path of the configuration file loaded from disk
func FileUsed() string {
	return viper.ConfigFileUsed()
}

const (
	TogglUsername  string = "toggl.username"
	TogglPassword  string = "toggl.password"
	TogglServerUrl string = "toggl.server.url"
	JiraServerUrl  string = "jira.server.url"
	JiraUsername   string = "jira.username"
	JiraPassword   string = "jira.password"
	JiraProjectKey string = "jira.project.key"
)

func Get(key string) string {
	return viper.GetString(key)
}

func Set(key string, value string) {
	viper.Set(key, value)
}

const jiraOverheadKeyPrefix = "jira.overhead"

// GetAllOverheadKeys returns all overhead keys from config, if any exist
func GetAllOverheadKeys() []string {
	overheadKeys := make([]string, 0)
	for _, key := range viper.AllKeys() {
		if keyName := strings.TrimPrefix(key, jiraOverheadKeyPrefix+"."); !strings.EqualFold(key, keyName) {
			overheadKeys = append(overheadKeys, keyName)
		}
	}
	return overheadKeys
}

// GetOverheadKey returns the specified overhead key from config, if any exists
func GetOverheadKey(key string) string {
	return viper.GetString(generateOverheadKeyFrom(key))
}

// SetOverheadKey accepts a new value for the specified overhead key to be stored in config
func SetOverheadKey(key string, value string) {
	viper.Set(generateOverheadKeyFrom(key), value)
}

func generateOverheadKeyFrom(key string) string {
	return fmt.Sprintf("%s.%s", jiraOverheadKeyPrefix, key)
}

// Reset clears all configuration loaded from disk (contents on disk are not removed)
func Reset() {
	viper.Reset()
}
