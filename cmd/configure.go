package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"syscall"
)

type InputController interface {
	RequestTextInput(string) (string, error)
	RequestPassword(string) (string, error)
}

type StdInReader struct{}

func (StdInReader) RequestTextInput(description string) (string, error) {
	fmt.Print(description)
	r := bufio.NewReader(os.Stdin)
	return r.ReadString('\n')
}

func (StdInReader) RequestPassword(description string) (string, error) {
	fmt.Print(description)
	bytes, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	return string(bytes), err
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create (or update) toggl-sync configuration",
	Long:  "Create (or update) the necessary configuration entries so all other toggl-sync commands work without issues",
	Run: func(cmd *cobra.Command, args []string) {
		err := configure(StdInReader{})
		if err != nil {
			log.Fatalf("Error configuring toggl-sync: %s", err)
		}
	},
}

func configure(reader InputController) error {
	err, _ := config.Init()
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading configuration file: %s", err))
	}

	err = updateConfiguration(reader)
	if err != nil {
		return errors.New(fmt.Sprintf("Error updating configuration: %s", err))
	}

	if err := config.Persist(); err != nil {
		return errors.New(fmt.Sprintf("Error saving configuration to file: %s", err))
	}

	return nil
}

func updateConfiguration(reader InputController) (err error) {
	config.SetTogglServerUrl("https://www.toggl.com/api/v8")
	err = saveSettingAs(reader, "Toggl username", config.GetTogglUsername, config.SetTogglUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(reader, "Toggl password", config.GetTogglPassword, config.SetTogglPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(reader, "Jira server url", config.GetJiraServerUrl, config.SetJiraServerUrl, false)
	if err != nil {
		return
	}
	err = saveSettingAs(reader, "Jira username", config.GetJiraUsername, config.SetJiraUsername, false)
	if err != nil {
		return
	}
	err = saveSettingAs(reader, "Jira password", config.GetJiraPassword, config.SetJiraPassword, true)
	if err != nil {
		return
	}
	err = saveSettingAs(reader, "Jira project key", config.GetJiraProjectKey, config.SetJiraProjectKey, false)
	if err != nil {
		return
	}
	return
}

func saveSettingAs(reader InputController, inputName string, getFn func() string, saveFn func(string), isPassword bool) error {
	existingValue := getFn()
	input, err := requestInput(reader, inputName, existingValue, isPassword)
	if err == nil && input != "" {
		saveFn(input)
	}
	return err
}

func requestInput(reader InputController, inputName string, existingValue string, isPassword bool) (string, error) {
	if isPassword {
		return requestPassword(reader, inputName, existingValue)
	} else {
		return requestTextInput(reader, inputName, existingValue)
	}
}

func requestTextInput(reader InputController, inputName string, existingValue string) (string, error) {
	var description string
	if existingValue != "" {
		description = fmt.Sprintf("%s (%s): ", inputName, existingValue)
	} else {
		description = fmt.Sprintf("%s: ", inputName)
	}

	input, err := reader.RequestTextInput(description)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading input: %s", err))
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func requestPassword(reader InputController, inputName string, existingValue string) (string, error) {
	var description string
	if existingValue != "" {
		description = fmt.Sprintf("%s (*****): ", inputName)
	} else {
		description = fmt.Sprintf("%s: ", inputName)
	}

	pwd, err := reader.RequestPassword(description)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading input: %s", err))
	}
	pwd = strings.TrimSpace(pwd)
	return pwd, nil
}
