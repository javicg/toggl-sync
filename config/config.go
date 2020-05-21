package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const togglUsernameKey = "toggl.username"
const togglPasswordKey = "toggl.password"
const jiraServerUrlKey = "jira.server.url"
const jiraUsernameKey = "jira.username"
const jiraPasswordKey = "jira.password"
const jiraProjectKeyKey = "jira.project.key"
const jiraOverheadKeyPrefix = "jira.overhead"

func Init() (err error, ok bool) {
	viper.SetConfigName("toggl-sync")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/usr/local/etc")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, false
		} else {
			return err, false
		}
	}

	return nil, true
}

func FileUsed() string {
	return viper.ConfigFileUsed()
}

func GetTogglUsername() string {
	return viper.GetString(togglUsernameKey)
}

func SetTogglUsername(username string) {
	viper.Set(togglUsernameKey, username)
}

func GetTogglPassword() string {
	return viper.GetString(togglPasswordKey)
}

func SetTogglPassword(password string) {
	viper.Set(togglPasswordKey, password)
}

func GetJiraServerUrl() string {
	return viper.GetString(jiraServerUrlKey)
}

func SetJiraServerUrl(serverUrl string) {
	viper.Set(jiraServerUrlKey, serverUrl)
}

func GetJiraUsername() string {
	return viper.GetString(jiraUsernameKey)
}

func SetJiraUsername(username string) {
	viper.Set(jiraUsernameKey, username)
}

func GetJiraPassword() string {
	return viper.GetString(jiraPasswordKey)
}

func SetJiraPassword(password string) {
	viper.Set(jiraPasswordKey, password)
}

func GetJiraProjectKey() string {
	return viper.GetString(jiraProjectKeyKey)
}

func SetJiraProjectKey(projectKey string) {
	viper.Set(jiraProjectKeyKey, projectKey)
}

func GetOverheadKey(key string) string {
	return viper.GetString(generateOverheadKeyFrom(key))
}

func SetOverheadKey(key string, mappedValue string) {
	viper.Set(generateOverheadKeyFrom(key), mappedValue)
}

func generateOverheadKeyFrom(key string) string {
	return fmt.Sprintf("%s.%s", jiraOverheadKeyPrefix, key)
}

func Persist() error {
	// Creating file beforehand as viper.WriteConfig fails otherwise
	err := createConfigFile()
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}

func createConfigFile() error {
	f, err := os.OpenFile("/usr/local/etc/toggl-sync.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	return f.Close()
}
