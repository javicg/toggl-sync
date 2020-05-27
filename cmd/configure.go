package cmd

import (
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create (or update) toggl-sync configuration",
	Long:  "Create (or update) the necessary configuration entries so all other toggl-sync commands work without issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := configure(stdInController{})
		if err != nil {
			log.Fatalf("Error configuring toggl-sync: %s", err)
		}
	},
}

func configure(inputCtrl inputController) error {
	err, _ := config.Init()
	if err != nil {
		return fmt.Errorf("error reading configuration file: %s", err)
	}

	err = updateConfiguration(inputCtrl)
	if err != nil {
		return fmt.Errorf("error updating configuration: %s", err)
	}

	if err := config.Persist(); err != nil {
		return fmt.Errorf("error saving configuration to file: %s", err)
	}

	return nil
}

func updateConfiguration(inputCtrl inputController) (err error) {
	config.SetTogglServerUrl("https://www.toggl.com/api/v8")
	err = saveSettingAs(inputCtrl, "Toggl username", config.GetTogglUsername, config.SetTogglUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Toggl password", config.GetTogglPassword, config.SetTogglPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira server url", config.GetJiraServerUrl, config.SetJiraServerUrl, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira username", config.GetJiraUsername, config.SetJiraUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira password", config.GetJiraPassword, config.SetJiraPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira project key", config.GetJiraProjectKey, config.SetJiraProjectKey, false)
	if err != nil {
		return
	}
	for _, key := range config.GetAllOverheadKeys() {
		getOverheadFn := func() string {
			return config.GetOverheadKey(key)
		}
		saveOverheadFn := func(value string) {
			config.SetOverheadKey(key, value)
		}
		if err = saveSettingAs(inputCtrl, fmt.Sprintf("Overhead - %s", key), getOverheadFn, saveOverheadFn, false); err != nil {
			return
		}
	}
	return
}

func saveSettingAs(inputCtrl inputController, inputName string, getFn func() string, saveFn func(string), isPassword bool) error {
	existingValue := getFn()
	input, err := requestInput(inputCtrl, inputName, existingValue, isPassword)
	if err == nil && input != "" {
		saveFn(input)
	}
	return err
}

func requestInput(inputCtrl inputController, inputName string, existingValue string, isPassword bool) (string, error) {
	if isPassword {
		return requestPassword(inputCtrl, inputName, existingValue)
	}

	return requestTextInput(inputCtrl, inputName, existingValue)
}

func requestTextInput(inputCtrl inputController, inputName string, existingValue string) (string, error) {
	var description string
	if existingValue != "" {
		description = fmt.Sprintf("%s (%s): ", inputName, existingValue)
	} else {
		description = fmt.Sprintf("%s: ", inputName)
	}

	input, err := inputCtrl.requestTextInput(description)
	if err != nil {
		return "", fmt.Errorf("error reading input: %s", err)
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func requestPassword(inputCtrl inputController, inputName string, existingValue string) (string, error) {
	var description string
	if existingValue != "" {
		description = fmt.Sprintf("%s (*****): ", inputName)
	} else {
		description = fmt.Sprintf("%s: ", inputName)
	}

	pwd, err := inputCtrl.requestPassword(description)
	if err != nil {
		return "", fmt.Errorf("error reading input: %s", err)
	}
	pwd = strings.TrimSpace(pwd)
	return pwd, nil
}
