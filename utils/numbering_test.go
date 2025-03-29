package utils

import (
	"NameTidy/testutils"
	"fmt"
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
		fileName     string
		digits       int
		hierarchical bool
		expected     string
	}{
		{"old_file.txt", 3, false, "001_old_file.txt"},
		{"old_file.txt", 1, false, "1_old_file.txt"},
		{"old_file.txt", 5, false, "00001_old_file.txt"},
		{"example.data.txt", 3, false, "001_example.data.txt"},
	}

	for _, tc := range testCases {
		oldFile := filepath.Join(dir, tc.fileName)
		file, err := os.Create(oldFile)
		if err != nil {
			t.Fatalf("ファイル作成に失敗: %v", err)
		}
		file.Close()

		// テスト: ファイル名に連番を追加
		newPath, err := AddNumbering(oldFile, tc.digits, tc.hierarchical)
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

func TestProcessDirectory(t *testing.T) {
	dir := "test_dir"
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// 通常のファイルリスト
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatalf("ファイル作成に失敗: %v", err)
		}
		f.Close()
	}

	// サブディレクトリ内のファイル作成
	subDir := filepath.Join(dir, "sub")
	os.Mkdir(subDir, 0755)
	subFiles := []string{"subfile1.txt", "subfile2.txt"}
	for _, file := range subFiles {
		f, err := os.Create(filepath.Join(subDir, file))
		if err != nil {
			t.Fatalf("サブディレクトリのファイル作成に失敗: %v", err)
		}
		f.Close()
	}

	// ディレクトリ内のファイルに連番を付ける
	err := ProcessDirectory(dir, 3, true) // 階層考慮
	if err != nil {
		t.Fatalf("ディレクトリの処理に失敗: %v", err)
	}

	// 各ファイルの名前が変更されたか確認
	for i, file := range files {
		expected := filepath.Join(dir, fmt.Sprintf("%03d_%s", i+1, file))
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("期待されるファイルが見つかりません: %s", expected)
		}
	}

	// サブディレクトリ内のファイル確認
	for i, file := range subFiles {
		expected := filepath.Join(subDir, fmt.Sprintf("%03d_%s", i+1, file))
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("期待されるファイルが見つかりません: %s", expected)
		}
	}
}

func TestProcessDirectory_EmptyDirectory(t *testing.T) {
	dir := "empty_dir"
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// ディレクトリ内のファイルに連番を付ける (空のディレクトリ)
	err := ProcessDirectory(dir, 3, false)
	if err != nil {
		t.Fatalf("空のディレクトリでエラーが発生しました: %v", err)
	}
}
