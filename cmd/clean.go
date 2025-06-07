package cmd

import (
	"NameTidy/internal/cleaner"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans up file names.",
    Run:   runWithCommonSetup("file name cleanup", cleaner.Clean),
}

func init() {
	cleanCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	cleanCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	cleanCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	cleanCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(cleanCmd)
}