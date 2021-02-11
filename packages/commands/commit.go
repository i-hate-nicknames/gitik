package commands

import (
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/commit"
	"github.com/spf13/cobra"
)

var messageP string

func init() {
	rootCmd.AddCommand(makeCommitCmd)
	makeCommitCmd.Flags().StringVarP(&messageP, "message", "m", "", "commit message")
}

var makeCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit changes to the repository",
	Long:  "write current tree with given message and store it separately",

	Run: func(cmd *cobra.Command, args []string) {
		data, err := commit.SaveCurrentTree(messageP)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	},
}
