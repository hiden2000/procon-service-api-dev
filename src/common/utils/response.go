package utils

import (
	"encoding/json"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
)

// SendJSONResponseは，クライアントにJSON形式のレスポンスを送信する．
// HTTPステータスコードと任意のデータをJSON形式でエンコードし，レスポンスボディに書き込む．
// JSONは整形された形式でエンコードされ，読みやすい形式でクライアントに提供される．
// エンコードに失敗した場合は，サーバーエラーをクライアントに送信する．
//
// パラメータ:
// - w http.ResponseWriter: レスポンスを書き込むためのHTTPレスポンスライター．
// - statusCode int: クライアントに送信するHTTPステータスコード．
// - data interface{}: クライアントに送信する任意のデータ．
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// JSONを整形してエンコード
	response, err := json.MarshalIndent(map[string]interface{}{
		"status":  statusCode,
		"result":  data,
		"message": nil,
	}, "", "    ") // 第二引数はプレフィックス（使用しない），第三引数はインデント
	if err != nil {
		// エンコード失敗時のエラーハンドリング
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// SendErrorResponseは，クライアントにエラー情報を含むJSON形式のレスポンスを送信する．
// HTTPステータスコードとエラーメッセージをJSON形式でエンコードし，レスポンスボディに書き込む．
// エラーメッセージはAPIErrorインターフェースを通じて統一された形式で提供される．
// JSONは整形された形式でエンコードされ，読みやすい形式でクライアントに提供される．
// エンコードに失敗した場合は，サーバーエラーをクライアントに送信する．
//
// パラメータ:
// - w http.ResponseWriter: レスポンスを書き込むためのHTTPレスポンスライター．
// - err error: クライアントに通知するエラー情報．
func SendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	// エラーをAPIErrorインターフェースに統一
	apiErr, statusCode := commonerrors.ErrorToAPIError(err), 0
	if apiErr == nil {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = apiErr.StatusCode
	}

	// JSONを整形してエンコード
	response, err := json.MarshalIndent(map[string]interface{}{
		"status": statusCode,
		"result": nil,
		"message": func() interface{} {
			if err != nil {
				return err.Error()
			} else {
				return nil
			}
		}(),
	}, "", "    ") // 第二引数はプレフィックス（使用しない），第三引数はインデント
	if err != nil {
		// エンコード失敗時のエラーハンドリング
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
