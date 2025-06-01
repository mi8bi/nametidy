package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Manage rename history",
}

// sub command: history clear
var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all rename history records",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		// For history clear, dirPath is not needed, so pass "" and checkDir=false
		db, err := handleCommonInitializations(verbose, "", false)
		if err != nil {
			utils.Error(err.Error(), nil)
			return
		}

		if err := cleaner.ClearHistory(db); err != nil {
			utils.Error("Failed to clear history", err)
		} else {
			utils.Info("History cleared")
		}
	},
}

func init() {
	historyClearCmd.Flags().BoolP("verbose", "v", false, "Show detailed logs")

	historyCmd.AddCommand(historyClearCmd)
	rootCmd.AddCommand(historyCmd)
}