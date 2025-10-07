package testutils

import (
	"os"
	"testing"
)

// setupTestEnvironment はテスト用のディレクトリを作成します
func SetupTestEnvironment(t *testing.T, dir string) {
	// テスト用ディレクトリを作成
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("テストディレクトリの作成に失敗: %v", err)
	}
}

// teardownTestEnvironment はテスト用のディレクトリを削除します
func TeardownTestEnvironment(t *testing.T, dir string) {
	// テスト用ディレクトリとその内容を削除
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("テストディレクトリの削除に失敗: %v", err)
	}
}
