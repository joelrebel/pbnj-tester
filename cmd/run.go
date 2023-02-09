package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joelrebel/pbnj-tester/internal"
	"github.com/spf13/cobra"
)

var (
	testsConfig string
	pbnjAddr    string
	bmcHost     string
	bmcUser     string
	bmcPass     string
	bmcPort     string
	logLevel    string
	testNames   []string
)

// runCmd represents the tester run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Test various PBnJ actions on one server",
	Run: func(cmd *cobra.Command, args []string) {
		tester := internal.NewTester(pbnjAddr, bmcHost, bmcUser, bmcPass, bmcPort, logLevel)
		tester.Run(cmd.Context(), testNames)

		b, err := json.MarshalIndent(tester.Results(), "", " ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	},
}

func init() {
	runCmd.PersistentFlags().StringVar(&pbnjAddr, "pbnj-addr", "localhost:50051", "PBnJ server port")
	runCmd.PersistentFlags().StringVar(&bmcHost, "bmc-ip", "", "BMC IP Address")
	runCmd.MarkPersistentFlagRequired("bmc-ip")

	runCmd.PersistentFlags().StringVar(&bmcUser, "bmc-user", "", "BMC username")
	runCmd.MarkPersistentFlagRequired("bmc-user")

	runCmd.PersistentFlags().StringVar(&bmcPass, "bmc-pass", "", "BMC password")
	runCmd.MarkPersistentFlagRequired("bmc-pass")

	runCmd.PersistentFlags().StringVar(&bmcPort, "bmc-port", "623", "BMC IPMI port")
	runCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level")

	runCmd.PersistentFlags().StringSliceVar(&testNames, "tests", []string{}, "test names to run")
	runCmd.MarkPersistentFlagRequired("tests")

	rootCmd.AddCommand(runCmd)
}
