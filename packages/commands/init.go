package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "start a new repository",
	Long:  "init creates a new repository in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initialized a new git repository")
	},
}
