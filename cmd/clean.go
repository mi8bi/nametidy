package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans up file names.",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		db, err := handleCommonInitializations(verbose, dirPath, true)
		if err != nil {
			utils.Error(err.Error(), nil) // Assuming utils.Error can take a string and nil error
			return
		}

		// --clean process
		utils.Info("Starting file name cleanup...")
		if err := cleaner.Clean(db, dirPath, dryRun); err != nil {
			utils.Error("File name cleanup failed", err)
			return
		}
		utils.Info("File name cleanup completed.")
	},
}

func init() {
	cleanCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	cleanCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	cleanCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	cleanCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(cleanCmd)
}