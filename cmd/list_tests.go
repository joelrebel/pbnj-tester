package cmd

import (
	"fmt"

	"github.com/joelrebel/pbnj-tester/internal"
	"github.com/spf13/cobra"
)

// actinsListCmd lists the PBnJ tests that the tester can run.
var testsListCmd = &cobra.Command{
	Use:   "list-tests",
	Short: "List the various PBnJ tests supported by tester",
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range internal.SupportedActions {
			fmt.Println(a)
		}
	},
}

func init() {
	rootCmd.AddCommand(testsListCmd)
}
