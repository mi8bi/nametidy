package cmd

import (
	"NameTidy/cleaner"
	"NameTidy/utils"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "リネームの取り消しを行います。",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// ロガーの初期化
		utils.InitLogger(false)

		// ディレクトリ存在チェック
		if !utils.IsDirectory(dirPath) {
			utils.Error("指定されたディレクトリが存在しません", nil)
			return
		}

		// --undo処理
		utils.Info("リネームの取り消しを開始します...")
		if err := cleaner.Undo(dirPath, dryRun); err != nil {
			utils.Error("リネームの取り消しに失敗しました", err)
			return
		}
		utils.Info("リネームの取り消しが完了しました。")
	},
}

func init() {
	undoCmd.Flags().StringP("path", "p", ".", "対象ディレクトリのパス")
	undoCmd.Flags().BoolP("dry-run", "d", false, "リネーム結果のみ表示")

	rootCmd.AddCommand(undoCmd)
}
