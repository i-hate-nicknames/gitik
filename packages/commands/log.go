package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/commit"
	"github.com/i-hate-nicknames/gitik/packages/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "get commit history",
	Long:  "get commit history ordered from newest to oldest",

	Run: func(cmd *cobra.Command, args []string) {
		var commitLog []commit.Commit
		var err error
		if len(args) > 0 {
			commitLog, err = commit.LogFrom(storage.OID(args[0]))
		} else {
			commitLog, err = commit.Log()
		}

		if errors.Is(err, commit.ErrNoHead) {
			fmt.Println("No commits found")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, commit := range commitLog {
			fmt.Printf("commit %s\n\n", commit.OID)
			fmt.Printf("    %s\n", commit.Message)
		}
	},
}
