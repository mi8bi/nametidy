package utils

import (
	"log"
)

var verboseMode = false

// InitLogger は詳細ログモードの設定を行う
func InitLogger(verbose bool) {
    verboseMode = verbose
}

// Info は情報ログを出力する
func Info(msg string) {
    if verboseMode {
        log.Println("[INFO] " + msg)
    }
}

// Warn は警告ログを出力する
func Warn(msg string) {
    log.Println("[WARN] " + msg)
}

// Error はエラーログを出力する
func Error(msg string, err error) {
    log.Printf("[ERROR] %s: %v\n", msg, err)
}
