package cmd

import (
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "toggl-sync [date]",
	Short: "Synchronize time entries to Jira",
	Long:  "Synchronize time entries to Jira using predefined project keys",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		syncDate, err := extractDateToSync(args)
		if err != nil {
			log.Fatalln(err)
		}
		readConfig()
		validateConfig()

		err = sync(stdInController{}, api.NewTogglApi(), api.NewJiraApi(), syncDate, dryRun)
		if err != nil {
			log.Fatalln(err)
		}

		if !dryRun {
			if err := config.Persist(); err != nil {
				log.Fatalln("Error saving configuration to file: ", err)
			}
		}
	},
}

func extractDateToSync(args []string) (syncDate string, err error) {
	if len(args) == 1 {
		return args[0], nil
	}
	if len(args) == 0 {
		return time.Now().Format("2006-01-02"), nil
	}
	return "", fmt.Errorf("requires 0 or 1 arg(s) but received %d", len(args))
}

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry-run toggl-sync (avoid side effects)")
}

// Execute runs the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func readConfig() {
	ok, err := config.Init()
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

func sync(inputCtrl inputController, togglApi api.TogglApi, jiraApi api.JiraApi, syncDate string, dryRun bool) error {
	err := printUserDetails(togglApi)
	if err != nil {
		return err
	}

	entries, err := getTimeEntriesForDate(togglApi, syncDate)
	if err != nil {
		return err
	}

	printSummary(syncDate, entries)

	ok, message := validateEntries(entries)
	if !ok {
		log.Print("Found issues during validation:")
		fmt.Print(message)
		log.Print("Please, correct the time entries above and try again.")
		return fmt.Errorf("validation failed")
	}

	if dryRun {
		log.Print("Logging work on Jira... SKIPPED! (dry-run)")
		return nil
	}

	logWorkOnJira(inputCtrl, togglApi, jiraApi, entries)
	return nil
}

func printUserDetails(togglApi api.TogglApi) error {
	log.Print("Fetching user details...")
	me, err := togglApi.GetMe()
	if err != nil {
		return fmt.Errorf("error fetching user details: %s", err)
	}

	log.Print("User details:")
	fmt.Printf("Name = %s, Email = %s\n", me.Data.Fullname, me.Data.Email)
	return nil
}

func getTimeEntriesForDate(togglApi api.TogglApi, dateStr string) ([]api.TimeEntry, error) {
	startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing input date: %s", err)
	}

	entries, err := togglApi.GetTimeEntries(startDate, startDate.AddDate(0, 0, 1))
	if err != nil {
		return nil, fmt.Errorf("error retrieving time entries: %s", err)
	}

	return entries, nil
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

func printSummary(syncDate string, entries []api.TimeEntry) {
	log.Printf("== Time Entries Summary (%s) ==", syncDate)
	for i := range entries {
		fmt.Printf("Entry: %s || Duration (s): %d\n", entries[i].Description, entries[i].Duration)
	}
}

func logWorkOnJira(inputCtrl inputController, togglApi api.TogglApi, jiraApi api.JiraApi, entries []api.TimeEntry) {
	log.Print("Logging work on Jira...")
	for _, entry := range entries {
		if isJiraTicket(entry) {
			logProjectWorkOnJira(jiraApi, entry)
		} else {
			logOverheadWorkOnJira(inputCtrl, togglApi, jiraApi, entry)
		}
	}
}

func logProjectWorkOnJira(jiraApi api.JiraApi, entry api.TimeEntry) {
	err := jiraApi.LogWork(entry.Description, time.Duration(entry.Duration)*time.Second)
	if err != nil {
		log.Printf("No time logged for [%s]; operation failed with an error: %s", entry.Description, err)
	} else {
		log.Printf("Successfully logged [%d]s for entry [%s]", entry.Duration, entry.Description)
	}
}

func logOverheadWorkOnJira(inputCtrl inputController, togglApi api.TogglApi, jiraApi api.JiraApi, entry api.TimeEntry) {
	project, err := togglApi.GetProjectById(entry.Pid)
	if err != nil {
		log.Printf("No time logged for [%s]; retrieving project information failed with an error: %s", entry.Description, err)
		return
	}

	if config.GetOverheadKey(project.Data.Name) == "" {
		err = requestOverheadKey(inputCtrl, entry, project)
		if err != nil {
			log.Printf("No time logged for [%s]; requesting project overhead key failed with an error: %s", entry.Description, err)
			return
		}
	}

	key := config.GetOverheadKey(project.Data.Name)
	err = jiraApi.LogWorkWithUserDescription(key, time.Duration(entry.Duration)*time.Second, entry.Description)
	if err != nil {
		log.Printf("No time logged for [%s] (project [%s]); operation failed with an error: %s", entry.Description, project.Data.Name, err)
	} else {
		log.Printf("Successfully logged [%d]s for entry [%s] (project [%s])", entry.Duration, entry.Description, project.Data.Name)
	}
}

func requestOverheadKey(inputCtrl inputController, entry api.TimeEntry, project *api.Project) error {
	description := fmt.Sprintf("No configuration found for entry [%s] (project [%s]). Which Jira ticket should be used for this type of work? -> ", entry.Description, project.Data.Name)
	input, err := inputCtrl.requestTextInput(description)
	if err != nil {
		return fmt.Errorf("error reading input: %s", err)
	}
	input = strings.TrimSpace(input)

	log.Printf("Saving configuration: entries for project [%s] will be tracked as [%s] from now on", project.Data.Name, input)
	config.SetOverheadKey(project.Data.Name, input)
	return nil
}
