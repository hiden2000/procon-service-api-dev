package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
	"strconv"

	"github.com/gorilla/mux"
)

// DecodeRequestBodyは，HTTPリクエストのボディをJSON形式から指定された構造体(dst)にデコードする．
// この関数は，リクエストボディの内容をAPIのエンドポイントで期待される形式にパースし，
// 指定された構造体のインスタンスにマッピングするために使用される．
//
// パラメータ:
// - r *http.Request: デコード対象のHTTPリクエスト．リクエストボディからJSONデータを読み取る．
// - dst interface{}: デコードしたデータを格納するための構造体のポインタ．この関数は，このポインタが指す構造体にデコード結果を格納する．
//
// 戻り値:
// - error: デコード処理中に発生したエラー．成功時はnil．
//
// この関数は，未知のフィールドが含まれている場合にエラーを発生させることで，APIが想定外のフィールドを受け取った場合の処理を安全に行うことを目的としている．
func DecodeRequestBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 未知のフィールドを拒否
	err := decoder.Decode(dst)
	if err != nil {
		return commonerrors.WrapRequestParsingError(err)
	}
	return nil
}

// GetIntVarFromRequestは，HTTPリクエストのURLパスから指定された変数名に対応する整数値を取得する．
// この関数は，URLパスパラメータから動的な値を抽出し，それを整数として利用するために使われる．
//
// パラメータ:
// - r *http.Request: 整数値を抽出する対象のHTTPリクエスト．
// - varName string: 抽出したい値の変数名．
//
// 戻り値:
// - int: 変数名に対応する整数値．
// - error: 変数が見つからない，または整数に変換できない場合のエラー．成功時はnil．
func GetIntVarFromRequest(r *http.Request, varName string) (int, error) {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars[varName])
	if err != nil {
		return 0, commonerrors.WrapRequestVariableError(varName, err) // strconv.Atoiが失敗した場合，エラーを返す
	}
	return value, nil
}

// GetStrVarFromRequestは，HTTPリクエストのURLパスから指定された変数名に対応する文字列値を取得する．
// この関数は，URLパスパラメータから動的な値を抽出し，それを文字列として利用するために使われる．
//
// パラメータ:
// - r *http.Request: 文字列値を抽出する対象のHTTPリクエスト．
// - varName string: 抽出したい値の変数名．
//
// 戻り値:
// - string: 変数名に対応する文字列値．
// - error: 変数が見つからない場合のエラー．成功時はnil．
func GetStrVarFromRequest(r *http.Request, varName string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[varName]
	if !ok {
		return "", commonerrors.WrapRequestVariableError(varName, fmt.Errorf("variable %s not found in request", varName))
	}
	return value, nil
}

// ParseProblemMetadataは，HTTPリクエストから問題メタデータを解析し，指定された構造体にデコードする．
// この関数は，問題の作成や更新に必要なメタデータをフォームデータから取得し，構造体にマッピングするために使用される．
//
// パラメータ:
// - r *http.Request: メタデータを含むHTTPリクエスト．
// - dst interface{}: デコードしたメタデータを格納する構造体のポインタ．
//
// 戻り値:
// - error: メタデータの解析またはデコードに失敗した場合のエラー．成功時はnil．
func ParseProblemMetadata(r *http.Request, dst interface{}) error {
	problemMetadata := r.FormValue("metadata")
	if problemMetadata == "" {
		return commonerrors.NewRequestParsingError("Metadata is required")
	}

	err := json.Unmarshal([]byte(problemMetadata), dst)
	if err != nil {
		return commonerrors.WrapRequestParsingError(err)
	}
	return nil
}

// ParseMultipartFormDataは，HTTPリクエストのマルチパートフォームデータをパースし，
// ファイルアップロードなどの処理を可能にする．この関数は，リクエストがマルチパートフォームデータを含む場合に使用される．
//
// パラメータ:
// - r *http.Request: パースする対象のHTTPリクエスト．
// - maxFileSize int64: アップロードを許可する最大ファイルサイズ．
//
// 戻り値:
// - error: マルチパートフォームデータのパースに失敗した場合のエラー．成功時はnil．
func ParseMultipartFormData(r *http.Request, maxFileSize int64) error {
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return commonerrors.NewValidationError("multipart_form", "Failed to parse request")
	}
	return nil
}
