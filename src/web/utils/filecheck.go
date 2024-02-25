package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	commonerrors "procon_web_service/src/common/errors"
)

// ValidateFilesは，マルチパートフォームデータに含まれる入力ファイルと出力ファイルの妥当性を検証する．
// 具体的には，ファイルの拡張子が.txtであること，入出力ファイル名が一致していること，およびファイル名が重複していないことを確認する．
// 検証に失敗した場合，対応するエラーを返す．
//
// パラメータ:
// - inputFiles []*multipart.FileHeader: 検証する入力ファイルのリスト．
// - outputFiles []*multipart.FileHeader: 検証する出力ファイルのリスト．
//
// 戻り値:
// - error: ファイル検証に失敗した場合のエラー．エラーは拡張子が.txtではないファイル，重複するファイル名，または入力ファイルに対応する出力ファイルが存在しない場合に発生する．
func ValidateFiles(inputFiles, outputFiles []*multipart.FileHeader) error {
	inputFileNames := make(map[string]bool)
	outputFileNames := make(map[string]bool)

	for _, file := range inputFiles {
		if filepath.Ext(file.Filename) != ".txt" {
			return commonerrors.NewFileValidationError("拡張子が .txt ではないファイルが存在します")
		}
		_, fileExists := inputFileNames[file.Filename]
		if fileExists {
			return commonerrors.NewFileValidationError("重複するファイル名が存在します")
		}
		inputFileNames[file.Filename] = true
	}

	for _, file := range outputFiles {
		if filepath.Ext(file.Filename) != ".txt" {
			return commonerrors.NewFileValidationError("拡張子が .txt ではないファイルが存在します")
		}
		_, fileExists := outputFileNames[file.Filename]
		if fileExists {
			return commonerrors.NewFileValidationError("重複するファイル名が存在します")
		}
		outputFileNames[file.Filename] = true
	}

	for inputFileName := range inputFileNames {
		if _, ok := outputFileNames[inputFileName]; !ok {
			return commonerrors.NewFileValidationError(fmt.Sprintf("input_file: %s に対応する output_file が存在しません", inputFileName))
		}
	}

	for outputFileName := range outputFileNames {
		if _, ok := inputFileNames[outputFileName]; !ok {
			return commonerrors.NewFileValidationError(fmt.Sprintf("output_file: %s に対応する input_file が存在しません", outputFileName))
		}
	}

	return nil
}
