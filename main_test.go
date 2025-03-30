package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const testDir = "test_data"

// buildExecutableでWindows環境用に修正
func buildExecutable(t *testing.T) string {
	exeName := "NameTidy"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", exeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ビルド失敗: %v\n出力: %s", err, string(output))
	}

	return exeName
}

// テスト用のディレクトリとファイルの準備
func setupTestEnvironment(t *testing.T) {
	os.Mkdir(testDir, 0755)
	files := []string{"IMG 2023 01 01.JPG", "_MyFile__.txt", "Special$$File!.docx", ".NameTidy_History"}
	for _, file := range files {
		path := filepath.Join(testDir, file)
		os.WriteFile(path, []byte("test content"), 0644)
	}
}

// テスト環境のクリーンアップ
func teardownTestEnvironment() {
	os.RemoveAll(testDir)
}

// TestClean - `--clean` のテスト
func TestClean(t *testing.T) {
	setupTestEnvironment(t)
	defer teardownTestEnvironment()

	exeName := buildExecutable(t)

	cmd := exec.Command("./" + exeName, "clean", "--path="+testDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("エラー: %v\n出力: %s", err, string(output))
	}

	expectedFiles := []string{"IMG_2023_01_01.JPG", "MyFile.txt", "Special_File.docx"}
	for _, file := range expectedFiles {
		if _, err := os.Stat(filepath.Join(testDir, file)); os.IsNotExist(err) {
			t.Errorf("期待されるファイルが見つかりません: %s", file)
		}
	}

	// `.NameTidy_History` は変更されていないことを確認
	if _, err := os.Stat(filepath.Join(testDir, ".NameTidy_History")); os.IsNotExist(err) {
		t.Errorf("履歴ファイルが見つかりません: .NameTidy_History")
	}
}

// TestDryRun - `--dry-run` のテスト
func TestDryRun(t *testing.T) {
	setupTestEnvironment(t)
	defer teardownTestEnvironment()

	exeName := buildExecutable(t)
	cmd := exec.Command("./" + exeName, "clean", "--path="+testDir, "--dry-run")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("エラー: %v\n出力: %s", err, string(output))
	}

	// 実際の出力からパスを削除
	actual := strings.Replace(string(output), testDir+"\\", "", -1)

	expectedOutput := []string{
		"[DRY-RUN] IMG 2023 01 01.JPG → IMG_2023_01_01.JPG",
		"[DRY-RUN] _MyFile__.txt → MyFile.txt",
		"[DRY-RUN] Special$$File!.docx → Special_File.docx",
	}

	for _, expected := range expectedOutput {
		if !strings.Contains(actual, expected) {
			t.Errorf("期待される出力が見つかりません: %s", expected)
		}
	}
}

// TestUndo - `--undo` のテスト
func TestUndo(t *testing.T) {
	setupTestEnvironment(t)
	defer teardownTestEnvironment()

	exeName := buildExecutable(t)

	// Step 1: `--clean` でリネーム実行
	cmd := exec.Command("./" + exeName, "clean", "--path="+testDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("エラー: %v\n出力: %s", err, string(output))
	}

	// Step 2: 履歴ファイルが存在するか確認
	historyFile := filepath.Join(testDir, ".NameTidy_History")
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Fatalf("履歴ファイルが存在しません: %s", historyFile)
	}

	// Step 3: `--undo` でリネーム取り消し
	cmd = exec.Command("./" + exeName, "undo", "--path="+testDir)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("エラー: %v\n出力: %s", err, string(output))
	}

	// 元のファイルが戻っていることを確認
	originalFiles := []string{"IMG 2023 01 01.JPG", "_MyFile__.txt", "Special$$File!.docx"}
	for _, file := range originalFiles {
		if _, err := os.Stat(filepath.Join(testDir, file)); os.IsNotExist(err) {
			t.Errorf("元のファイルが見つかりません: %s", file)
		}
	}
}

// TestUndoDryRun - `--undo --dry-run` のテスト
func TestUndoDryRun(t *testing.T) {
	setupTestEnvironment(t)
	defer teardownTestEnvironment()

	exeName := buildExecutable(t)

	// Step 1: `--clean` でリネーム実行
	cmd := exec.Command("./" + exeName, "clean", "--path="+testDir)
	cmd.CombinedOutput()

	// Step 2: `--undo --dry-run` で取り消しのシミュレーション
	cmd = exec.Command("./" + exeName, "undo", "--path="+testDir, "--dry-run")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("エラー: %v\n出力: %s", err, string(output))
	}

	// 実際の出力からパスを削除
	actual := strings.Replace(string(output), testDir+"\\", "", -1)

	expectedOutput := []string{
		"[DRY-RUN] IMG_2023_01_01.JPG → IMG 2023 01 01.JPG",
		"[DRY-RUN] MyFile.txt → _MyFile__.txt",
		"[DRY-RUN] Special_File.docx → Special$$File!.docx",
	}

	for _, expected := range expectedOutput {
		if !strings.Contains(actual, expected) {
			t.Errorf("期待される出力が見つかりません: %s", expected)
		}
	}
}

// TestInvalidDirectory - エラーハンドリングのテスト
func TestInvalidDirectory(t *testing.T) {
	exeName := buildExecutable(t)
	cmd := exec.Command("./" + exeName, "clean", "--path=invalid_dir")
	output, _ := cmd.CombinedOutput()

	if !strings.Contains(string(output), "The specified directory does not exist") {
		t.Errorf("存在しないディレクトリのエラーメッセージが正しくありません。出力: %s", string(output))
	}
}
