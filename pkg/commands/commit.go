package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/pkg/commit"
	"github.com/i-hate-nicknames/gitik/pkg/storage"
	"github.com/spf13/cobra"
)

var messageP string

func init() {
	rootCmd.AddCommand(makeCommitCmd)
	makeCommitCmd.Flags().StringVarP(&messageP, "message", "m", "", "commit message")
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(checkoutCmd)
}

var makeCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit changes to the repository",
	Long:  "write current tree with given message and store it separately",

	Run: func(cmd *cobra.Command, args []string) {
		treeOID, err := commit.SaveCurrentTree(messageP)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(treeOID[:]))
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "get commit history",
	Long:  "get commit history ordered from newest to oldest",

	Run: func(cmd *cobra.Command, args []string) {
		var commitLog []commit.Commit
		var err error
		if len(args) > 0 {
			commitOID, err := storage.MakeOID([]byte(args[0]))
			if err != nil {
				log.Fatal(err)
			}
			commitLog, err = commit.LogFrom(commitOID)
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

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "check out given commit, resetting working tree to it",
	Long:  "set working tree to the tree of the commit and update HEAD",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("Expecting commit hash")
		}
		commitOID, err := storage.MakeOID([]byte(args[0]))
		if err != nil {
			log.Fatalf(err.Error())
		}
		c, err := commit.GetCommit(commitOID)
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = c.Checkout()
		if err != nil {
			log.Fatalf(err.Error())
		}
	},
}
