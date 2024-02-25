package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"procon_web_service/src/common/config"
	"strings"
)

// SaveCodeToFile - 提出されたソースコードを一時ファイルに保存
func SaveCodeToFile(languageID int, code string) (codeFilePath string, langConfig config.LanguageConfig, cleanupFunc func() error, err error) {
	langConfig, ok := config.GetLanguageConfigByID(languageID)
	if !ok {
		return "", langConfig, nil, fmt.Errorf("unsupported language ID: %d", languageID)
	}

	codeDir, err := ioutil.TempDir("/tmp", "codeDir_")
	if err != nil {
		return "", langConfig, nil, fmt.Errorf("failed to create temporary directory for code: %v", err)
	}

	cleanupFunc = func() error {
		return os.RemoveAll(codeDir)
	}

	fileExt := getFileExtension(langConfig.LanguageID)
	codeFile, err := os.Create(filepath.Join(codeDir, fmt.Sprintf("solution.%s", fileExt)))
	if err != nil {
		return "", langConfig, cleanupFunc, fmt.Errorf("failed to create code file: %v", err)
	}
	defer codeFile.Close()

	// コードをファイルに書き込み(\n　は改行として認識)
	if _, err := codeFile.WriteString(strings.Replace(code, "\\n", "\n", -1)); err != nil {
		return "", langConfig, cleanupFunc, fmt.Errorf("failed to write code to file: %v", err)
	}

	return codeFile.Name(), langConfig, cleanupFunc, nil
}

// CreateTempFile - 一時的なディレクトリ内部に一時的なファイルを作成
func CreateTempFile() (tempFilePath string, cleanupFunc func() error, err error) {
	tempDirPath, err := ioutil.TempDir("/tmp", "judge_")
	if err != nil {
		return "", nil, err
	}

	cleanupFunc = func() error {
		return os.RemoveAll(tempDirPath)
	}

	tempFilePath = filepath.Join(tempDirPath, "tempfile.txt")
	if err := ioutil.WriteFile(tempFilePath, []byte(""), 0666); err != nil {
		cleanupFunc()
		return "", nil, err
	}

	return tempFilePath, cleanupFunc, nil
}

// getFileExtension - 言語IDに基づいてファイル拡張子を計算
func getFileExtension(languageID int) string {
	switch languageID {
	case 1: // Python
		return "py"
	case 2: // C++
		return "cpp"
	case 3: // Go
		return "go"
	case 4: // Java
		return "java"
	case 5: // Rust
		return "rs"
	default:
		return "txt"
	}
}
