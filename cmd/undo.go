package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the most recent rename operation.",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Initialize logger
		utils.InitLogger(verbose)

		// Check if directory exists
		if !utils.IsDirectory(dirPath) {
			utils.Error("The specified directory does not exist", nil)
			return
		}

		// --undo process
		utils.Info("Starting to undo the rename operation...")
		if err := cleaner.Undo(dirPath, dryRun); err != nil {
			utils.Error("Failed to undo the rename operation", err)
			return
		}
		utils.Info("Undoing the rename operation is complete.")
	},
}

func init() {
	undoCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	undoCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	undoCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	undoCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(undoCmd)
}
