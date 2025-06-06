/*
Copyright © 2025 mi8bi <mi8biiiii@gmail.com>
*/
package cmd

import (
	"NameTidy/internal/cleaner"
	"github.com/spf13/cobra"
)

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redoes the most recent rename operation.",
    Run:   runWithCommonSetup("redo the rename operation", cleaner.Redo),
}

func init() {
	redoCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	redoCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	redoCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	redoCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(redoCmd)
}
