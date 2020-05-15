package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
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
			log.Fatalln("No configuration file exists! Please, run 'configure' to create a new configuration file")
		} else {
			log.Fatalf("Unable to read configuration: %s", err)
		}
	}

	log.Printf("Configuration read from: %s", viper.ConfigFileUsed())
}

func validateConfig() {
	isValid :=
		valueExists("TOGGL_USERNAME") &&
			valueExists("TOGGL_PASSWORD") &&
			valueExists("JIRA_SERVER_URL") &&
			valueExists("JIRA_USERNAME") &&
			valueExists("JIRA_PASSWORD")

	if !isValid {
		log.Fatalln("Configuration file is invalid! Please, run 'configure' to create a new configuration file")
	}
}

func valueExists(configName string) bool {
	if value := viper.Get(configName); value == nil {
		log.Printf("%s not specified!", configName)
		return false
	}
	return true
}

func sync() {
	togglApi := api.NewTogglApi()

	log.Print("Fetching user details...")
	me, err := togglApi.GetMe()
	if err != nil {
		log.Fatalf("Error fetching user details: %s", err)
	}

	log.Printf("User details: Name = %s, Email = %s", me.Data.Fullname, me.Data.Email)

	fmt.Print("Introduce a date to fetch time entries (e.g. 2020-05-08) -> ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
	input = strings.Replace(input, "\n", "", -1)

	startDate, err := time.Parse(time.RFC3339, input+"T00:00:00Z")
	if err != nil {
		log.Fatalf("Error parsing input date: %s", err)
	}

	entries, err := togglApi.GetTimeEntries(startDate, startDate.AddDate(0, 0, 1))
	if err != nil {
		log.Fatalf("Error retrieving time entries: %s", err)
		return
	}

	log.Print("== Time Entries Summary ==")
	for i := range entries {
		log.Printf("Entry: %s || Duration (s): %d", entries[i].Description, entries[i].Duration)
	}
}
