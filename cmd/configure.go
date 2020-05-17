package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create (or update) toggl-sync configuration",
	Long:  "Create (or update) the necessary configuration entries so all other toggl-sync commands work without issues",
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func configure() {
	err, _ := config.Init()
	if err != nil {
		log.Fatalln("Error reading configuration file: ", err)
	}

	updateConfiguration()
	if err := config.Persist(); err != nil {
		log.Fatalln("Error saving configuration to file: ", err)
	}
}

func updateConfiguration() {
	saveSettingAs("Toggl username", config.GetTogglUsername, config.SetTogglUsername, false)
	saveSettingAs("Toggl password", config.GetTogglPassword, config.SetTogglPassword, true)
	saveSettingAs("Jira server url", config.GetJiraServerUrl, config.SetJiraServerUrl, false)
	saveSettingAs("Jira username", config.GetJiraUsername, config.SetJiraUsername, false)
	saveSettingAs("Jira password", config.GetJiraPassword, config.SetJiraPassword, true)
	saveSettingAs("Jira project key", config.GetJiraProjectKey, config.SetJiraProjectKey, false)
}

func saveSettingAs(inputName string, getFn func() string, saveFn func(string), isPassword bool) {
	existingValue := getFn()
	input := requestInput(inputName, existingValue, isPassword)
	if input != "" {
		saveFn(input)
	}
}

func requestInput(inputName string, existingValue string, isPassword bool) string {
	if isPassword {
		return requestPassword(inputName, existingValue)
	} else {
		return requestTextInput(inputName, existingValue)
	}
}

func requestTextInput(inputName string, existingValue string) string {
	if existingValue != "" {
		fmt.Printf("%s (%s): ", inputName, existingValue)
	} else {
		fmt.Printf("%s: ", inputName)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("Error reading input:", err)
	}
	input = strings.TrimSpace(input)
	return input
}

func requestPassword(inputName string, existingValue string) string {
	if existingValue != "" {
		fmt.Printf("%s (*****): ", inputName)
	} else {
		fmt.Printf("%s: ", inputName)
	}

	bytePwd, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		log.Fatalln("Error reading input:", err)
	}
	pwd := string(bytePwd)
	pwd = strings.TrimSpace(pwd)
	return pwd
}
