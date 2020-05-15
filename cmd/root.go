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
	Use:   "toggl-sync",
	Short: "Synchornize time entries to Jira",
	Long:  "Synchronize time entries to Jira using predefined project keys",
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
			config.GetJiraPassword() != "" &&
			config.GetJiraProjectKey() != ""

	if !isValid {
		log.Fatalln("Configuration file is invalid! Please, run 'configure' to create a new configuration file")
	}
}

func sync() {
	togglApi := api.NewTogglApi()
	jiraApi := api.NewJiraApi()

	printUserDetails(togglApi)

	trackingDate := requestTrackingDate()
	entries := getTimeEntriesForDate(togglApi, trackingDate)

	printSummary(entries)
	logWorkOnJira(togglApi, jiraApi, entries)
	if err := config.Persist(); err != nil {
		log.Fatalln("Error saving configuration to file: ", err)
	}
}

func printUserDetails(togglApi *api.TogglApi) {
	log.Print("Fetching user details...")
	me, err := togglApi.GetMe()
	if err != nil {
		log.Fatalf("Error fetching user details: %s", err)
	}

	log.Printf("User details: Name = %s, Email = %s", me.Data.Fullname, me.Data.Email)
}

func requestTrackingDate() string {
	fmt.Print("Introduce a date to fetch time entries (e.g. 2020-05-08) -> ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
	input = strings.Replace(input, "\n", "", -1)
	return input
}

func getTimeEntriesForDate(togglApi *api.TogglApi, dateStr string) []api.TimeEntry {
	startDate, err := time.Parse(time.RFC3339, dateStr+"T00:00:00Z")
	if err != nil {
		log.Fatalf("Error parsing input date: %s", err)
	}

	entries, err := togglApi.GetTimeEntries(startDate, startDate.AddDate(0, 0, 1))
	if err != nil {
		log.Fatalf("Error retrieving time entries: %s", err)
	}

	return entries
}

func printSummary(entries []api.TimeEntry) {
	log.Print("== Time Entries Summary ==")
	for i := range entries {
		log.Printf("Entry: %s || Duration (s): %d", entries[i].Description, entries[i].Duration)
	}
}

func logWorkOnJira(togglApi *api.TogglApi, jiraApi *api.JiraApi, entries []api.TimeEntry) {
	log.Print("Logging work on Jira...")
	for _, entry := range entries {
		if strings.HasPrefix(entry.Description, config.GetJiraProjectKey()) {
			logProjectWorkOnJira(jiraApi, entry)
		} else {
			logOverheadWorkOnJira(togglApi, jiraApi, entry)
		}
	}
}

func logProjectWorkOnJira(jiraApi *api.JiraApi, entry api.TimeEntry) {
	err := jiraApi.LogWork(entry.Description, time.Duration(entry.Duration)*time.Second)
	if err != nil {
		log.Printf("No time logged for [%s]; operation failed with an error: %s", entry.Description, err)
	} else {
		log.Printf("Successfully logged [%d]s for entry [%s]", entry.Duration, entry.Description)
	}
}

func logOverheadWorkOnJira(togglApi *api.TogglApi, jiraApi *api.JiraApi, entry api.TimeEntry) {
	project := togglApi.GetProjectById(entry.Pid)
	if key := config.GetOverheadKey(project.Data.Name); key == "" {
		requestOverheadKey(entry, project)
	}

	key := config.GetOverheadKey(project.Data.Name)
	err := jiraApi.LogWork(key, time.Duration(entry.Duration)*time.Second)
	if err != nil {
		log.Printf("No time logged for [%s] (project [%s]); operation failed with an error: %s", entry.Description, project.Data.Name, err)
	} else {
		log.Printf("Successfully logged [%d]s for entry [%s] (project [%s])", entry.Duration, entry.Description, project.Data.Name)
	}
}

func requestOverheadKey(entry api.TimeEntry, project api.ProjectData) {
	fmt.Printf("No configuration found for entry [%s] (project [%s]). Which Jira ticket should be used for this type of work? -> ", entry.Description, project.Data.Name)
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
	input = strings.Replace(input, "\n", "", -1)

	log.Printf("Saving configuration: entries for project [%s] will be tracked as [%s] from now on", project.Data.Name, input)
	config.SetOverheadKey(project.Data.Name, input)
}
