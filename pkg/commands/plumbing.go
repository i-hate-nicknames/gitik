package commands

import (
	"fmt"
	"log"

	"github.com/i-hate-nicknames/gitik/pkg/plumbing"
	"github.com/i-hate-nicknames/gitik/pkg/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readTreeCmd)
	rootCmd.AddCommand(writeTreeCmd)
	rootCmd.AddCommand(hashObjCmd)
	rootCmd.AddCommand(catFileCmd)
}

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "retrieve an object from the index",
	Long:  "find an object by its object id in the index and print its contents to stdout",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		obj, err := storage.GetObject(storage.OID(args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(obj.Data))
	},
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
