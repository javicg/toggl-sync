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

// Available configuration keys
const (
	TogglUsername  string = "toggl.username"
	TogglPassword  string = "toggl.password"
	TogglServerURL string = "toggl.server.url"
	JiraServerURL  string = "jira.server.url"
	JiraUsername   string = "jira.username"
	JiraPassword   string = "jira.password"
	JiraProjectKey string = "jira.project.key"
)

// Get returns the current value of the key in the config map (if any exists)
func Get(key string) string {
	return viper.GetString(key)
}

// GetSlice returns the current values associated with the key in the config map (if any exist)
func GetSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// Set overrides the value of the key in the config map
func Set(key string, value interface{}) {
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
