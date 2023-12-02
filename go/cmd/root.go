package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aoc23",
	Short: "Provides cli access to aoc solutions.",
	Long: `Use this to either run the solutions, or generate a new folder and files with correct day specific folder
name and package names everywhere.

Usage:
$ aoc23 run
$ aoc23 run 11
$ aoc23 generate 21`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
