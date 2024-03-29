package commands

import (
	"fmt"

	"github.com/i-hate-nicknames/gitik/pkg/storage"
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
		path := storage.Init()
		fmt.Printf("Initialized empty gitik repository in %s\n", path)
	},
}
