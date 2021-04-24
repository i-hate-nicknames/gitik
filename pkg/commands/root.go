package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitik",
	Short: "gitik is a small tiny reimplementation of git",
	Long:  "gitik is a small tiny reimplementation of git, serving educational purposes",
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
