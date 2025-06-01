/*
Copyright © 2025 mi8bi <mi8biiiii@gmail.com>
*/
package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redoes the most recent rename operation.",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		db, err := handleCommonInitializations(verbose, dirPath, true)
		if err != nil {
			utils.Error(err.Error(), nil)
			return
		}

		// --redo process
		utils.Info("Starting to redo the rename operation...")
		if err := cleaner.Redo(db, dirPath, dryRun); err != nil {
			utils.Error("Failed to redo the rename operation", err)
			return
		}
		utils.Info("Undoing the rename operation is complete.")
	},
}

func init() {
	redoCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	redoCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	redoCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	redoCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(redoCmd)
}
