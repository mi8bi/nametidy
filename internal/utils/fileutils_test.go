package utils

import (
	"NameTidy/testutils"
	"os"
	"path/filepath"
	"testing"
)

func TestIsDirectory(t *testing.T) {
	// テスト用のディレクトリをセットアップ
	dir := "test_dir"
	testFile := "test_file.txt"
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// ファイルを作成
	file, err := os.Create(filepath.Join(dir, testFile))
	if err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}
	file.Close()

	// ディレクトリの場合
	if !IsDirectory(dir) {
		t.Errorf("expected %s to be a directory", dir)
	}

	// ファイルの場合
	if IsDirectory(filepath.Join(dir, testFile)) {
		t.Errorf("expected %s to be a file, not a directory", filepath.Join(dir, testFile))
	}
}

func TestCleanFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My_File___.txt", "My_File.txt"},
		{"Special$$File!.docx", "Special_File.docx"},
		{"IMG 2023 01 01.JPG", "IMG_2023_01_01.JPG"},
		{"_MyFile__.txt", "MyFile.txt"},
		{"__My__File__.txt", "My_File.txt"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := CleanFileName(test.input)
			if output != test.expected {
				t.Errorf("expected %s, got %s", test.expected, output)
			}
		})
	}
}

func TestRenameFileDryRun(t *testing.T) {
	// テスト用のディレクトリをセットアップ
	dir := "test_dir"
	oldFile := filepath.Join(dir, "old_file.txt")
	newFile := filepath.Join(dir, "new_file.txt")
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// 新しいファイルを作成
	file, err := os.Create(oldFile)
	if err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}
	file.Close()

	// dryRunがtrueの場合、ファイルはリネームされない
	err = RenameFile(oldFile, newFile, true)
	if err != nil {
		t.Errorf("dryRunモードでエラー: %v", err)
	}

	// dryRunがtrueの場合、実際にはファイルがリネームされない
	if FileExists(newFile) {
		t.Errorf("dryRunモードでファイルがリネームされました: %s", newFile)
	}
}

func TestRenameFileRealRename(t *testing.T) {
	// テスト用のディレクトリとファイルを作成
	dir := "rename_real_test_dir"
	oldFile := filepath.Join(dir, "old_file.txt")
	newFile := filepath.Join(dir, "new_file.txt")
	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// 新しいファイルを作成
	file, err := os.Create(oldFile)
	if err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}
    file.Close()

	// dryRunがfalseの場合、実際にファイルをリネーム
	err = RenameFile(oldFile, newFile, false)
	if err != nil {
		t.Fatalf("ファイルリネームに失敗: %v", err)
	}
}

func TestFileExists(t *testing.T) {
	// テスト用のファイルパス
	dir := "test_dir"
	test_file := "test_file.txt"
	filePath := filepath.Join(dir, test_file)

	testutils.SetupTestEnvironment(t, dir)
	defer testutils.TeardownTestEnvironment(t, dir)

	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}
	file.Close()

	// ファイルが存在する場合
	if !FileExists(filePath) {
		t.Errorf("expected %s to exist", filePath)
	}

	// 存在しないファイル
	if FileExists(filepath.Join(dir, "non_existent_file.txt")) {
		t.Errorf("expected file to not exist")
	}
}
