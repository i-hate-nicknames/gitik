package commands

import (
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hashObjCmd)
}

var hashObjCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "hash and store given file in the index",
	Long:  "calculate sha1 hash of the contents of the given file and store it under this hash",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		oid, err := plumbing.WriteFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(oid)
	},
}
