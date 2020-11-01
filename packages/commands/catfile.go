package commands

import (
	"bytes"
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/packages/data"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(catFileCmd)
}

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "retrieve an object from the index",
	Long:  "find an object by its object id in the index and print its contents to stdout",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		data, err := data.GetObject(args[0], data.TypeBlob)
		if err != nil {
			log.Fatal(err)
		}
		buf := bytes.NewBuffer(data)
		fmt.Printf(buf.String())
	},
}
