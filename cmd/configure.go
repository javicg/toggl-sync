package cmd

import (
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"strings"
)

// NewConfigureCmd creates a new Cobra Command that helps configuring the application
func NewConfigureCmd(configManager config.Manager, inputCtrl inputController) *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Create (or update) toggl-sync configuration",
		Long:  "Create (or update) the necessary configuration entries so all other toggl-sync commands work without issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := configure(configManager, inputCtrl)
			return err
		},
	}
}

func configure(configManager config.Manager, inputCtrl inputController) error {
	_, err := configManager.Init()
	if err != nil {
		return fmt.Errorf("error reading configuration file: %s", err)
	}

	err = updateConfiguration(inputCtrl)
	if err != nil {
		return fmt.Errorf("error updating configuration: %s", err)
	}

	if err := configManager.Persist(); err != nil {
		return fmt.Errorf("error saving configuration to file: %s", err)
	}

	return nil
}

func updateConfiguration(inputCtrl inputController) (err error) {
	config.Set(config.TogglServerURL, "https://api.track.toggl.com/api/v8")
	err = saveSettingAs(inputCtrl, "Toggl username", config.TogglUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Toggl password", config.TogglPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira server url", config.JiraServerURL, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira username", config.JiraUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira password", config.JiraPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(inputCtrl, "Jira project key", config.JiraProjectKey, false)
	if err != nil {
		return
	}
	for _, key := range config.GetAllOverheadKeys() {
		if err = saveOverheadSettingAs(inputCtrl, fmt.Sprintf("Overhead - %s", key), key); err != nil {
			return
		}
	}
	return
}

func saveSettingAs(inputCtrl inputController, inputName string, key string, isPassword bool) error {
	existingValue := config.Get(key)
	input, err := requestInput(inputCtrl, inputName, existingValue, isPassword)
	if err == nil && input != "" {
		config.Set(key, input)
	}
	return err
}

func saveOverheadSettingAs(inputCtrl inputController, inputName string, key string) error {
	existingValue := config.GetOverheadKey(key)
	input, err := requestTextInput(inputCtrl, inputName, existingValue)
	if err == nil && input != "" {
		config.SetOverheadKey(key, input)
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
