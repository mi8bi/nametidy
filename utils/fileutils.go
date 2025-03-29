package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsDirectory は指定パスがディレクトリかを確認する
func IsDirectory(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return info.IsDir()
}

// CleanFileName 修正
func CleanFileName(fileName string) string {
	// 拡張子を取り除く
	ext := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(ext)]

	// アルファベット、数字、ドット以外の文字を _ に置き換える
	reClean := regexp.MustCompile(`[^\w\d.]`)
	baseName = reClean.ReplaceAllString(baseName, "_")

	// 連続するアンダースコアを1つにする
	reUnderscore := regexp.MustCompile(`_+`)
	baseName = reUnderscore.ReplaceAllString(baseName, "_")

	// 先頭と末尾のアンダースコアを取り除く
	baseName = strings.Trim(baseName, "_")

	// 拡張子を元に戻す
	return baseName + ext
}



// RenameFile はファイルをリネームする
func RenameFile(oldPath, newPath string, dryRun bool) error {
    if dryRun {
        fmt.Printf("[DRY-RUN] %s → %s\n", oldPath, newPath)
        return nil
    }

    err := os.Rename(oldPath, newPath)
    if err != nil {
        return fmt.Errorf("リネームに失敗しました: %v", err)
    }

    fmt.Printf("Renamed: %s → %s\n", filepath.Base(oldPath), filepath.Base(newPath))
    return nil
}

// FileExists は指定したパスのファイルが存在するか確認する
func FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}