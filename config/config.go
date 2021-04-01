package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const togglUsernameKey = "toggl.username"
const togglPasswordKey = "toggl.password"
const togglServerUrlKey = "toggl.server.url"
const jiraServerUrlKey = "jira.server.url"
const jiraUsernameKey = "jira.username"
const jiraPasswordKey = "jira.password"
const jiraProjectKeyKey = "jira.project.key"
const jiraOverheadKeyPrefix = "jira.overhead"

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

// FileUsed returns the absolute path of the configuration file loaded from disk
func FileUsed() string {
	return viper.ConfigFileUsed()
}

// GetTogglUsername returns the Toggl username from config, if any exists
func GetTogglUsername() string {
	return viper.GetString(togglUsernameKey)
}

// SetTogglUsername accepts a new Toggl username to be stored in config
func SetTogglUsername(username string) {
	viper.Set(togglUsernameKey, username)
}

// GetTogglPassword returns the Toggl password from config, if any exists
func GetTogglPassword() string {
	return viper.GetString(togglPasswordKey)
}

// SetTogglPassword accepts a new Toggl password to be stored in config
func SetTogglPassword(password string) {
	viper.Set(togglPasswordKey, password)
}

// GetTogglServerUrl returns the Toggl server url from config, if any exists (used as a base url for API calls)
func GetTogglServerUrl() string {
	return viper.GetString(togglServerUrlKey)
}

// SetTogglServerUrl accepts a new Toggl server url to be stored in config
func SetTogglServerUrl(serverUrl string) {
	viper.Set(togglServerUrlKey, serverUrl)
}

// GetJiraServerUrl returns the Jira server url from config, if any exists (used as a base url for API calls)
func GetJiraServerUrl() string {
	return viper.GetString(jiraServerUrlKey)
}

// SetJiraServerUrl accepts a new Jira server url to be stored in config
func SetJiraServerUrl(serverUrl string) {
	viper.Set(jiraServerUrlKey, serverUrl)
}

// GetJiraUsername returns the Jira username from config, if any exists
func GetJiraUsername() string {
	return viper.GetString(jiraUsernameKey)
}

// SetJiraUsername accepts a new Jira username to be stored in config
func SetJiraUsername(username string) {
	viper.Set(jiraUsernameKey, username)
}

// GetJiraPassword returns the Jira password from config, if any exists
func GetJiraPassword() string {
	return viper.GetString(jiraPasswordKey)
}

// SetJiraPassword accepts a new Jira password to be stored in config
func SetJiraPassword(password string) {
	viper.Set(jiraPasswordKey, password)
}

// GetJiraProjectKey returns the Jira project key from config, if any exists
func GetJiraProjectKey() string {
	return viper.GetString(jiraProjectKeyKey)
}

// SetJiraProjectKey accepts a new Jira project key to be stored in config
func SetJiraProjectKey(projectKey string) {
	viper.Set(jiraProjectKeyKey, projectKey)
}

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
func SetOverheadKey(key string, mappedValue string) {
	viper.Set(generateOverheadKeyFrom(key), mappedValue)
}

func generateOverheadKeyFrom(key string) string {
	return fmt.Sprintf("%s.%s", jiraOverheadKeyPrefix, key)
}

// Reset clears all configuration loaded from disk (contents on disk are not removed)
func Reset() {
	viper.Reset()
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
