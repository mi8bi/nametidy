package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// AddNumbering ファイル名に連番を付与する
func AddNumbering(path string, digits int, index int) (string, error) {
	dir, file := filepath.Split(path)
	newName := GenerateNumberedName(file, digits, index)

	// ファイルの新しいパスを作成
	newPath := filepath.Join(dir, newName)
	return newPath, nil
}

// GenerateNumberedName ファイル名に番号を付ける
func GenerateNumberedName(baseName string, digits int, index int) string {
	indexStr := fmt.Sprintf("%0*d", digits, index)
	ext := filepath.Ext(baseName)
	fileName := baseName[:len(baseName)-len(ext)]
	return fmt.Sprintf("%s_%s%s", indexStr, fileName, ext)
}

// ProcessDirectory ディレクトリ内のファイルに連番を付ける
func ProcessDirectory(dirPath string, digits int, hierarchical bool) error {
	counts := make(map[string]int)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var dirKey string
			if hierarchical {
				dirKey = filepath.Dir(path)
			} else {
				dirKey = "global"
			}

			// カウントのインクリメント
			counts[dirKey]++
			count := counts[dirKey]

			// ファイルの連番を付ける
			newPath, err := AddNumbering(path, digits, count)
			if err != nil {
				return err
			}

			// ファイルをリネーム
			if err := os.Rename(path, newPath); err != nil {
				return fmt.Errorf("ファイルのリネームに失敗しました: %v", err)
			}
			fmt.Printf("Renamed: %s → %s\n", path, newPath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("ディレクトリの処理中にエラーが発生しました: %v", err)
	}
	return nil
}