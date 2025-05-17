package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var numberCmd = &cobra.Command{
	Use:   "number",
	Short: "Adds sequence numbers to file names.",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		numbered, _ := cmd.Flags().GetInt("numbered")
		hierarchical, _ := cmd.Flags().GetBool("hierarchical")
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

		// --numbered process
		utils.Info("Starting to add sequence numbers to file names...")
		if err := cleaner.NumberFiles(db, dirPath, numbered, hierarchical, dryRun); err != nil {
			utils.Error("Failed to add sequence numbers to file names", err)
			return
		}
		utils.Info("Sequence number addition to file names completed.")
	},
}

func init() {
	numberCmd.Flags().StringP("path", "p", "", "Path to the target directory")
	numberCmd.Flags().BoolP("dry-run", "d", false, "Show rename results only")
	numberCmd.Flags().IntP("numbered", "n", 3, "Add sequence numbers to file names")
	numberCmd.Flags().BoolP("hierarchical", "H", false, "Add sequence numbers based on directory structure")
	numberCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")
	numberCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(numberCmd)
}