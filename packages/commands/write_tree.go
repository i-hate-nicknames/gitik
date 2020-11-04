package commands

import (
	"log"

	"github.com/i-hate-nicknames/gitik/packages/base"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write given tree to the object database",
	Long:  "find an object by its object id in the index and print its contents to stdout",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		_, err := base.WriteTree(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}
