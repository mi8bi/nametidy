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

		// Initialize logger
		utils.InitLogger(verbose)

		// Check if directory exists
		if !utils.IsDirectory(dirPath) {
			utils.Error("The specified directory does not exist", nil)
			return
		}

		// Initialize DB
		db, err := cleaner.GetDB()
		if err != nil {
			utils.Error("Failed to open DB", err)
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