package cleaner

import (
	"NameTidy/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const HISTORY_FILE = ".NameTidy_History"

// Clean は指定ディレクトリ内のファイル名をクリーンアップする（再帰的にサブディレクトリも処理）
func Clean(dirPath string, dryRun bool) error {
	entries := make(map[string]string)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 履歴ファイルは無視
		if info.IsDir() || filepath.Base(path) == HISTORY_FILE {
			return nil
		}

		oldName := info.Name()
		newName := utils.CleanFileName(oldName)

		if oldName != newName {
			newPath := filepath.Join(filepath.Dir(path), newName)
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", path, newPath)
			} else {
				if err := os.Rename(path, newPath); err != nil {
					utils.Error("リネーム失敗", err)
					return err
				}
				fmt.Printf("Renamed: %s → %s\n", path, newPath)
				entries[path] = newPath
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

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

	historyFilePath := filepath.Join(dirPath, HISTORY_FILE)
	fmt.Printf("履歴ファイルパス: %s\n", historyFilePath)

	return os.WriteFile(historyFilePath, data, 0644)
}

// NumberFiles ファイル名へのナンバリングを行う
func NumberFiles(dirPath string, digits int, hierarchical bool, dryRun bool) error {
	// ディレクトリごとのカウントを管理
	counts := make(map[string]int)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 履歴ファイルは無視
		if info.IsDir() || filepath.Base(path) == HISTORY_FILE {
			return nil
		}

		var dirKey string
		if hierarchical {
			// ディレクトリごとにカウントをリセット
			dirKey = filepath.Dir(path)
		} else {
			// すべてのファイルで通し番号
			dirKey = "global"
		}

		// カウントのインクリメント
		counts[dirKey]++
		count := counts[dirKey]

		// 新しいファイル名の生成
		newPath, err := utils.AddNumbering(path, digits, count)
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
		return nil
	})
}

// Undo は直前のリネーム操作を取り消す
func Undo(dirPath string, dryRun bool) error {
	history, err := loadHistory(dirPath)
	if err != nil {
		return err
	}

	for oldPath, newPath := range history {
		if oldPath == HISTORY_FILE || newPath == HISTORY_FILE {
			continue // 履歴ファイルは無視
		}

		if utils.FileExists(newPath) {
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", newPath, oldPath)
			} else {
				err := utils.RenameFile(newPath, oldPath, false)
				if err != nil {
					utils.Error("リネームの取り消し失敗", err)
				}
			}
		}
	}

	if !dryRun {
		return os.Remove(filepath.Join(dirPath, HISTORY_FILE))
	}

	return nil
}

// loadHistory はリネーム履歴を読み込む
func loadHistory(dirPath string) (map[string]string, error) {
	data, err := os.ReadFile(filepath.Join(dirPath, HISTORY_FILE))
	if err != nil {
		return nil, errors.New("リネーム履歴が見つかりません")
	}

	history := make(map[string]string)
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, errors.New("履歴データの読み取りに失敗しました")
	}

	return history, nil
}
