package main

import (
	"github.com/javicg/toggl-sync/api"
	"github.com/javicg/toggl-sync/cmd"
	"github.com/javicg/toggl-sync/config"
	"log"
)

func main() {
	configManager := &config.ViperConfigManager{}
	inputCtrl := cmd.StdInController{}

	rootCmd := cmd.NewRootCmd(configManager, inputCtrl, api.NewTogglApi(), api.NewJiraApi())
	rootCmd.AddCommand(cmd.NewConfigureCmd(configManager, inputCtrl))
	rootCmd.AddCommand(cmd.NewVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
