package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "toggl-sync",
	Run: func(cmd *cobra.Command, args []string) {
		readConfig()
		validateConfig()
		sync()
	},
}

func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No configuration file exists! Please, run 'configure' to create a new configuration file")
			os.Exit(0)
		} else {
			fmt.Println("Unable to read configuration: ", err)
			os.Exit(1)
		}
	}

	fmt.Println("Configuration read from: ", viper.ConfigFileUsed())
}

func validateConfig() {
	isValid :=
		valueExists("TOGGL_USERNAME") &&
			valueExists("TOGGL_PASSWORD") &&
			valueExists("JIRA_SERVER_URL") &&
			valueExists("JIRA_USERNAME") &&
			valueExists("JIRA_PASSWORD")

	if !isValid {
		fmt.Println("Configuration file is invalid! Please, run 'configure' to create a new configuration file")
		os.Exit(1)
	}
}

func valueExists(configName string) bool {
	if value := viper.Get(configName); value == nil {
		fmt.Printf(fmt.Sprintf("%s not specified!\n", configName))
		return false
	}
	return true
}

func sync() {
	togglApi := api.NewTogglApi()

	fmt.Println("Fetching user details...")
	me, err := togglApi.GetMe()
	if err != nil {
		return
	}
	fmt.Printf("User details: Name = %s, Email = %s\n", me.Data.Fullname, me.Data.Email)

	fmt.Print("Introduce a date to fetch time entries (e.g. 2020-05-08) -> ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)

	startDate, err := time.Parse(time.RFC3339, input+"T00:00:00Z")
	if err != nil {
		fmt.Println("Error parsing input date:", err)
		return
	}

	entries, err := togglApi.GetTimeEntries(startDate, startDate.AddDate(0, 0, 1))
	if err != nil {
		return
	}

	fmt.Println("== Time Entries Summary ==")
	for i := range entries {
		fmt.Printf("Entry: %s || Duration (s): %d\n", entries[i].Description, entries[i].Duration)
	}
}
