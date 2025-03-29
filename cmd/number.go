package cmd

import (
	"NameTidy/cleaner"
	"NameTidy/utils"

	"github.com/spf13/cobra"
)

var numberCmd = &cobra.Command{
	Use:   "number",
	Short: "ファイル名に連番を追加します。",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		numbered, _ := cmd.Flags().GetInt("numbered")

		// ロガーの初期化
		utils.InitLogger(false)

		// ディレクトリ存在チェック
		if !utils.IsDirectory(dirPath) {
			utils.Error("指定されたディレクトリが存在しません", nil)
			return
		}

		// --numbered処理
		utils.Info("ファイル名への連番追加を開始します...")
		if err := cleaner.NumberFiles(dirPath, numbered, false, dryRun); err != nil {
			utils.Error("ファイル名への連番追加に失敗しました", err)
			return
		}
		utils.Info("ファイル名への連番追加が完了しました。")
	},
}

func init() {
	numberCmd.Flags().StringP("path", "p", ".", "対象ディレクトリのパス")
	numberCmd.Flags().BoolP("dry-run", "d", false, "リネーム結果のみ表示")
	numberCmd.Flags().IntP("numbered", "n", 0, "ファイル名に連番を付ける")

	rootCmd.AddCommand(numberCmd)
}