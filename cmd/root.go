package cmd

import (
	"bufio"
	"fmt"
	"github.com/javicg/toggl-sync/api"
	"github.com/spf13/cobra"
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
		sync()
	},
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