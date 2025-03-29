package cleaner

import (
	"NameTidy/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const historyFile = ".NameTidy_History"

// Clean は指定ディレクトリ内のファイル名をクリーンアップする（再帰的にサブディレクトリも処理）
func Clean(dirPath string, dryRun bool) error {
	// Walkで再帰的にファイルとディレクトリを走査
	entries := make(map[string]string) // 変更前後のファイル名を保持するマップ

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}

		// ファイル名のクリーンアップ
		oldName := info.Name()
		newName := utils.CleanFileName(oldName)

		// 名前が変更される場合
		if oldName != newName {
			newPath := filepath.Join(filepath.Dir(path), newName)
			// dry-runの場合、実際にはリネームしない
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", path, newPath)
			} else {
				if err := os.Rename(path, newPath); err != nil {
					utils.Error("リネーム失敗", err)
					return err
				}
				fmt.Printf("Renamed: %s → %s\n", path, newPath)
				// 実際にリネームした場合、履歴に追加
				entries[oldName] = newName
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// dryRunがfalseの場合のみ履歴を保存
	if !dryRun {
		if err := saveHistory(dirPath, entries); err != nil {
			utils.Error("履歴の保存に失敗しました", err)
			return err
		}
	}

	return nil
}

// saveHistory はリネーム履歴を保存
func saveHistory(dirPath string, history map[string]string) error {
    data, err := json.Marshal(history)
    if err != nil {
        return err
    }

    // 履歴ファイルの保存パスを確認
    historyFilePath := filepath.Join(dirPath, historyFile)
    fmt.Printf("履歴ファイルパス: %s\n", historyFilePath)  // 履歴ファイルパスをログ出力

    return os.WriteFile(historyFilePath, data, 0644)
}

// NumberFiles ファイル名へのナンバリングを行う
func NumberFiles(dirPath string, digits int, hierarchical bool, dryRun bool) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			newPath, err := utils.AddNumbering(path, digits, hierarchical)
			if err != nil {
				return err
			}
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", path, newPath)
			} else {
				if err := os.Rename(path, newPath); err != nil {
					return fmt.Errorf("ファイルのリネームに失敗しました: %v", err)
				}
				fmt.Printf("Renamed: %s → %s\n", path, newPath)
			}
		}
		return nil
	})
}

// Undo は直前のリネーム操作を取り消す
func Undo(dirPath string) error {
    history, err := loadHistory(dirPath)
    if err != nil {
        return err
    }

    for newName, oldName := range history {
        oldPath := filepath.Join(dirPath, oldName)
        newPath := filepath.Join(dirPath, newName)

        if utils.FileExists(oldPath) {
            err := utils.RenameFile(oldPath, newPath, false)
            if err != nil {
                utils.Error("リネームの取り消し失敗", err)
            }
        }
    }

    // 履歴ファイルの削除
    return os.Remove(filepath.Join(dirPath, historyFile))
}

// loadHistory はリネーム履歴を読み込む
func loadHistory(dirPath string) (map[string]string, error) {
    data, err := os.ReadFile(filepath.Join(dirPath, historyFile))
    if err != nil {
        return nil, errors.New("リネーム履歴が見つかりません")
    }

    history := make(map[string]string)
    if err := json.Unmarshal(data, &history); err != nil {
        return nil, errors.New("履歴データの読み取りに失敗しました")
    }

    return history, nil
}
