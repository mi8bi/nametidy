package utils

import (
	"nametidy/testutils"
	"os"
	"path/filepath"
	"testing"
)

func TestAddNumbering(t *testing.T) {
	// テスト用のディレクトリをセットアップ
	dir := "test_dir"
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// 通常のファイル作成
	testCases := []struct {
		fileName string
		digits   int
		index    int
		expected string
	}{
		{"old_file.txt", 3, 1, "001_old_file.txt"},
		{"old_file.txt", 1, 1, "1_old_file.txt"},
		{"old_file.txt", 5, 1, "00001_old_file.txt"},
		{"example.data.txt", 3, 1, "001_example.data.txt"},
	}

	for _, tc := range testCases {
		oldFile := filepath.Join(dir, tc.fileName)
		file, err := os.Create(oldFile)
		if err != nil {
			t.Fatalf("ファイル作成に失敗: %v", err)
		}
		file.Close()

		// テスト: ファイル名に連番を追加
		newPath, err := AddNumbering(oldFile, tc.digits, tc.index)
		if err != nil {
			t.Fatalf("ファイル名に連番を付ける処理に失敗: %v", err)
		}

		// 期待されるファイル名
		expectedPath := filepath.Join(dir, tc.expected)
		if newPath != expectedPath {
			t.Errorf("期待されるパス: %s, 実際のパス: %s", expectedPath, newPath)
		}
	}
}
