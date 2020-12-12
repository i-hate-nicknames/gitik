package commands

import (
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write given tree to the object database",
	Long:  "Given directory path, generate a tree object out of it and write it to the object database",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		oid, err := plumbing.WriteTree(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(oid)
	},
}
