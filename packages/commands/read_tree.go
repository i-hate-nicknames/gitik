package commands

import (
	"log"

	"github.com/i-hate-nicknames/gitik/packages/plumbing"
	"github.com/i-hate-nicknames/gitik/packages/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readTreeCmd)
}

var readTreeCmd = &cobra.Command{
	Use:   "read-tree",
	Short: "read tree from the object database",
	Long:  "Given oid of a tree, read it from the database recursively and write all files/directories into current directory",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		err := plumbing.ReadTree(storage.OID(args[0]))
		if err != nil {
			log.Fatal(err)
		}

	},
}
