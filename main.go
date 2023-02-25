package main

import (
	"log"

	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/cmd"
	"github.com/javicg/toggl-sync/config"
)

func main() {
	configManager := &config.ViperConfigManager{}
	inputCtrl := cmd.StdInController{}

	rootCmd := cmd.NewRootCmd(configManager, inputCtrl, api.NewTogglAPI(), api.NewJiraAPI())
	rootCmd.AddCommand(cmd.NewConfigureCmd(configManager, inputCtrl))
	rootCmd.AddCommand(cmd.NewVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
