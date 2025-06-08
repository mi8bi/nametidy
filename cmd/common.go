package cmd

import (
    "nametidy/internal/cleaner"
    "nametidy/internal/utils"
    "github.com/spf13/cobra"
    "gorm.io/gorm"
)

type operationFunc func(db *gorm.DB, dirPath string, dryRun bool) error

func runWithCommonSetup(opName string, op operationFunc) func(cmd *cobra.Command, args []string) {
    return func(cmd *cobra.Command, args []string) {
        dirPath, _ := cmd.Flags().GetString("path")
        dryRun, _ := cmd.Flags().GetBool("dry-run")
        verbose, _ := cmd.Flags().GetBool("verbose")

        utils.InitLogger(verbose)

        if !utils.IsDirectory(dirPath) {
            utils.Error("The specified directory does not exist", nil)
            return
        }

        db, err := cleaner.GetDB()
        if err != nil {
            utils.Error("Failed to open DB", err)
            return
        }

        utils.Info("Starting " + opName + "...")
        if err := op(db, dirPath, dryRun); err != nil {
            utils.Error(opName+" failed", err)
            return
        }
        utils.Info(opName + " completed.")
    }
}