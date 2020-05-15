package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		configure()
	},
}

func configure() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("No configuration file exists. Creating a new one...")
			createConfigFile("config.yaml")
		} else {
			log.Fatalln("Unable to read configuration: ", err)
		}
	}

	updateConfiguration()
	if err := viper.WriteConfig(); err != nil {
		log.Fatalln("Error saving configuration to file: ", err)
	}
}

func createConfigFile(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Error creating configuration file: ", err)
	}
	if err = f.Close(); err != nil {
		log.Fatalln("Error closing file: ", err)
	}
}

func updateConfiguration() {
	saveSettingAs("Toggl username", "TOGGL_USERNAME")
	savePasswordAs("Toggl password", "TOGGL_PASSWORD")
	saveSettingAs("Jira server url", "JIRA_SERVER_URL")
	saveSettingAs("Jira username", "JIRA_USERNAME")
	savePasswordAs("Jira password", "JIRA_PASSWORD")
}

func savePasswordAs(inputName string, key string) {
	value := viper.GetString(key)
	input := requestUserInput(inputName, value, true)
	if input != "" {
		viper.Set(key, input)
	}
}

func saveSettingAs(inputName string, key string) {
	value := viper.GetString(key)
	input := requestUserInput(inputName, value, false)
	if input != "" {
		viper.Set(key, input)
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
