package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-yaml/yaml"
	"github.com/joelrebel/pbnj-tester/internal"
	"github.com/spf13/cobra"
)

// runMultipleCmd represents the tester command to run PBnJ actions on multiple hosts
var runMultipleCmd = &cobra.Command{
	Use:   "run-multiple",
	Short: "Test various PBnJ actions on servers listed in the tests-config",
	Run: func(cmd *cobra.Command, args []string) {
		runMultiple(cmd.Context())
	},
}

func runMultiple(ctx context.Context) {
	b, err := os.ReadFile(testsConfig)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &internal.Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		log.Fatal(err)
	}

	testers := []*internal.Tester{}

	if len(cfg.Tests) == 0 {
		log.Fatal("no tests defined in configuration")
	}

	if len(cfg.Servers) == 0 {
		log.Fatal("no servers defined in configuration")
	}

	resultStore := internal.NewTestResultStore()

	wg := &sync.WaitGroup{}

	for _, server := range cfg.Servers {
		server := server
		tester := internal.NewTester(
			pbnjAddr,
			server.BmcHost,
			server.BmcUser,
			server.BmcPass,
			server.IpmiPort,
			logLevel,
		)

		testers = append(testers, tester)

		wg.Add(1)
		go func() {
			defer wg.Done()
			tester.Run(ctx, cfg.Tests)

			result := internal.DeviceResult{
				Vendor:  server.Vendor,
				Model:   server.Model,
				Name:    server.Name,
				BMCIP:   server.BmcHost,
				Results: tester.Results(),
			}

			resultStore.Save(result)
		}()
	}

	log.Println("waiting for tests to complete...")

	wg.Wait()

	results := resultStore.Read()

	b, err = json.MarshalIndent(results, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func init() {
	runMultipleCmd.PersistentFlags().StringVar(&pbnjAddr, "pbnj-addr", "localhost:50051", "PBnJ server port")
	runMultipleCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level")

	runMultipleCmd.PersistentFlags().StringVar(&testsConfig, "tests-config", "", "YAML file with test configuration")
	runMultipleCmd.MarkPersistentFlagRequired("tests-config")

	rootCmd.AddCommand(runMultipleCmd)
}
