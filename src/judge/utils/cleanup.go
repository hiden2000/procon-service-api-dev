package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StartCleanupScheduler - 定期的に/tmpディレクトリをクリーンアップするスケジューラを開始
func StartCleanupScheduler(interval time.Duration, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			cleanupTmpFiles("/tmp", maxAge)
		}
	}()
}

// cleanupTmpFiles - 指定されたディレクトリ内の古いファイルを削除
func cleanupTmpFiles(directory string, maxAge time.Duration) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Printf("Failed to list directory %s: %v", directory, err)
		return
	}

	now := time.Now()
	for _, file := range files {
		if now.Sub(file.ModTime()) > maxAge {
			os.Remove(filepath.Join(directory, file.Name()))
			log.Printf("Deleted old file: %s", file.Name())
		}
	}
}
