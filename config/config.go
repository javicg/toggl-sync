package config

import (
	"github.com/spf13/viper"
	"os"
)

func Init() (err error, ok bool) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, false
		} else {
			return err, false
		}
	}

	return nil, true
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

func GetTogglUsername() string {
	return viper.GetString("TOGGL_USERNAME")
}

func SetTogglUsername(username string) {
	viper.Set("TOGGL_USERNAME", username)
}

func GetTogglPassword() string {
	return viper.GetString("TOGGL_PASSWORD")
}

func SetTogglPassword(password string) {
	viper.Set("TOGGL_PASSWORD", password)
}

func GetJiraServerUrl() string {
	return viper.GetString("JIRA_SERVER_URL")
}

func SetJiraServerUrl(serverUrl string) {
	viper.Set("JIRA_SERVER_URL", serverUrl)
}

func GetJiraUsername() string {
	return viper.GetString("JIRA_USERNAME")
}

func SetJiraUsername(username string) {
	viper.Set("JIRA_USERNAME", username)
}

func GetJiraPassword() string {
	return viper.GetString("JIRA_PASSWORD")
}

func SetJiraPassword(password string) {
	viper.Set("JIRA_PASSWORD", password)
}

func GetJiraProjectKey() string {
	return viper.GetString("JIRA_PROJECT_KEY")
}

func SetJiraProjectKey(projectKey string) {
	viper.Set("JIRA_PROJECT_KEY", projectKey)
}

func Persist() error {
	// Creating file beforehand as viper.WriteConfig fails otherwise
	err := createConfigFile("config.yaml")
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}

func createConfigFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	return f.Close()
}
