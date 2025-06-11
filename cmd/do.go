package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "mark task as done",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will mark a task as done")
	},
}
