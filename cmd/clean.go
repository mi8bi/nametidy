package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "ファイル名をクリーンアップします。",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// ロガーの初期化
		utils.InitLogger(verbose)

		// ディレクトリ存在チェック
		if !utils.IsDirectory(dirPath) {
			utils.Error("指定されたディレクトリが存在しません", nil)
			return
		}

		// --clean処理
		utils.Info("ファイル名のクリーニングを開始します...")
		if err := cleaner.Clean(dirPath, dryRun); err != nil {
			utils.Error("ファイル名のクリーニングに失敗しました", err)
			return
		}
		utils.Info("ファイル名のクリーニングが完了しました。")
	},
}

func init() {
	cleanCmd.Flags().StringP("path", "p", "", "対象ディレクトリのパス")
	cleanCmd.Flags().BoolP("dry-run", "d", false, "リネーム結果のみ表示")
	cleanCmd.Flags().BoolP("verbose", "v", false, "詳細なログを表示")
	cleanCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(cleanCmd)
}