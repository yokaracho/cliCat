package root

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cliCat",
	Short: "A command-line tool to fetch cat breeds information.",
	Long: `cliCat is a CLI tool for fetching cat breeds information from an API,
grouping them by country, sorting breed names by length, and saving the result to a JSON file.`,
	Run: runCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		getLogger.Error("Error executing root command:", err)
		os.Exit(1)
	}
}
