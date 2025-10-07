package cmd

import (
	"nametidy/internal/cleaner"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the most recent rename operation.",
	Run:   runWithCommonSetup("undo the rename operation", cleaner.Undo),
}

func init() {
	undoCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	undoCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	undoCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	undoCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(undoCmd)
}
