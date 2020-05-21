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

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "toggl-sync",
	Short: "Synchronize time entries to Jira",
	Long:  "Synchronize time entries to Jira using predefined project keys",
	Run: func(cmd *cobra.Command, args []string) {
		readConfig()
		validateConfig()
		sync(dryRun)
	},
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry-run toggl-sync (avoid side effects)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func readConfig() {
	err, ok := config.Init()
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err)
	}

	if !ok {
		log.Fatalln("No configuration file exists! Please, run 'configure' to create a new configuration file")
	}

	log.Printf("Configuration read from: %s", config.FileUsed())
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

func sync(dryRun bool) {
	togglApi := api.NewTogglApi()

	printUserDetails(togglApi)

	trackingDate := requestTrackingDate()
	entries := getTimeEntriesForDate(togglApi, trackingDate)

	printSummary(entries)

	ok, message := validateEntries(entries)
	if !ok {
		log.Print("Found issues during validation:")
		fmt.Print(message)
		log.Fatal("Please, correct the time entries above and try again.")
	}

	if dryRun {
		log.Print("Logging work on Jira... SKIPPED! (dry-run)")
	} else {
		jiraApi := api.NewJiraApi()
		logWorkOnJira(togglApi, jiraApi, entries)
		if err := config.Persist(); err != nil {
			log.Fatalln("Error saving configuration to file: ", err)
		}
	}
}

func printUserDetails(togglApi *api.TogglApi) {
	log.Print("Fetching user details...")
	me, err := togglApi.GetMe()
	if err != nil {
		log.Fatalf("Error fetching user details: %s", err)
	}

	log.Print("User details:")
	fmt.Printf("Name = %s, Email = %s\n", me.Data.Fullname, me.Data.Email)
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

func validateEntries(entries []api.TimeEntry) (ok bool, message string) {
	log.Print("Validating time entries...")
	ok, message = true, ""
	for _, entry := range entries {
		entryOk, entryMessage := validateEntry(entry)
		ok = ok && entryOk
		message = message + entryMessage
	}
	return
}

func validateEntry(entry api.TimeEntry) (ok bool, message string) {
	if entry.Description == "" {
		return false, "Found entry without a description. All entries must contain a description.\n"
	} else if !isJiraTicket(entry) && entry.Pid == 0 {
		return false, fmt.Sprintf("Entry [%s] does not seem to be a Jira ticket and doesn't have a Toggl project assigned.\n", entry.Description)
	} else {
		return true, ""
	}
}

func isJiraTicket(entry api.TimeEntry) bool {
	return strings.HasPrefix(entry.Description, config.GetJiraProjectKey())
}

func printSummary(entries []api.TimeEntry) {
	log.Print("== Time Entries Summary ==")
	for i := range entries {
		fmt.Printf("Entry: %s || Duration (s): %d\n", entries[i].Description, entries[i].Duration)
	}
}

func logWorkOnJira(togglApi *api.TogglApi, jiraApi *api.JiraApi, entries []api.TimeEntry) {
	log.Print("Logging work on Jira...")
	for _, entry := range entries {
		if isJiraTicket(entry) {
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
	err := jiraApi.LogWorkWithUserDescription(key, entry.Description, time.Duration(entry.Duration)*time.Second)
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
