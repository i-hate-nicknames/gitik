package commands

import (
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/data"
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
		oid, err := data.HashObject(args[0], data.TypeBlob)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(oid)
	},
}
