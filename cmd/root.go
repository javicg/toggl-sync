package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
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
	err, ok := config.Init()
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err)
	}

	if !ok {
		log.Fatalln("No configuration file exists! Please, run 'configure' to create a new configuration file")
	}

	log.Printf("Configuration read from: %s", config.ConfigFileUsed())
}

func validateConfig() {
	isValid :=
		config.GetTogglUsername() != "" &&
			config.GetTogglPassword() != "" &&
			config.GetJiraServerUrl() != "" &&
			config.GetJiraUsername() != "" &&
			config.GetJiraPassword() != ""

	if !isValid {
		log.Fatalln("Configuration file is invalid! Please, run 'configure' to create a new configuration file")
	}
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
