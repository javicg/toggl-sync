package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use: "configure",
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
	value := getFn()
	input := requestUserInput(inputName, value, isPassword)
	if input != "" {
		saveFn(input)
	}
}

func requestUserInput(inputName string, previousValue string, isPassword bool) (input string) {
	if previousValue != "" && !isPassword {
		fmt.Printf("%s (%s): ", inputName, previousValue)
	} else if previousValue != "" && isPassword {
		fmt.Printf("%s (*****): ", inputName)
	} else {
		fmt.Printf("%s: ", inputName)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("Error reading input:", err)
	}
	input = strings.Replace(input, "\n", "", -1)
	return
}
