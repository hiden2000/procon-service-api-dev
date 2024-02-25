package async

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"procon_web_service/src/common/models"
	"procon_web_service/src/web/config"
	"procon_web_service/src/web/database"

	"github.com/gorilla/websocket"
)

var (
	judgeURL = config.NewAPIEndpointsConfig().JudgeServerURL
)

// JudgeSolutionAsyncは，WebSocketを使用して解答の非同期判定を行い，結果をクライアントに通知する関数である．
// この関数は，提出された解答をジャッジサーバーへ送信し，判定結果を取得した後，その結果をWebSocketを介してクライアントに送信する．
// 判定プロセス中に発生したエラーは，WebSocketを通じてクライアントにエラーメッセージとして送信される．
// 解答の判定結果は，ジャッジサーバーからのレスポンスとして受け取り，models.ResultDetail構造体にデシリアライズされる．
// 最後に，判定結果をデータベースに保存し，WebSocketを使用してクライアントに判定結果を送信する．
// この関数は，WebSocket通信を介してユーザーにリアルタイムのフィードバックを提供するための非同期処理の一部として機能する．
//
// パラメータ:
// - ctx context.Context: 操作の実行に使用されるコンテキスト．
// - db *sql.DB: データベース接続へのポインタ．
// - solution models.Solution: 判定する解答．
// - conn *websocket.Conn: クライアントとのWebSocket接続．
//
// 注意:
// - この関数は，ジャッジサーバーへのリクエスト送信，レスポンスの処理，結果のクライアントへの送信を行う．
// - 判定結果のデータベースへの保存は，本番環境でのみ実行されるべきであり，開発やテスト環境では異なる扱いが必要になる場合がある．
func JudgeSolutionAsync(ctx context.Context, db *sql.DB, solution models.Solution, conn *websocket.Conn) {
	solutionBytes, err := json.Marshal(solution)
	if err != nil {
		SendError(conn, "Failed to marshal solution: "+err.Error())
		return
	}

	// ジャッジサーバーへのリクエストを準備
	respBytes, err := sendRequestToJudgeServer(ctx, judgeURL, solutionBytes)
	if err != nil {
		errMsg := "Failed to send request to judge server: " + err.Error()
		// タイムアウトによるエラーメッセージの調整
		if errors.Is(err, context.Canceled) {
			errMsg = "Judge server request timed out"
		}
		SendError(conn, errMsg)
		return
	}

	// ジャッジサーバーからのレスポンスを取得
	var response struct {
		Status  int                 `json:"status"`
		Result  models.ResultDetail `json:"result"`
		Message string              `json:"message"`
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		SendError(conn, "Failed to unmarshal judge server response: "+err.Error())
		return
	}

	resultDetail := response.Result

	if err := database.CreateResultDetail(db, solution.SolutionID, &resultDetail); err != nil {
		SendError(conn, "Failed to save result detail: "+err.Error())
		return
	}

	// WebSocketを通じて結果をクライアントに送信
	if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
		SendError(conn, "Failed to send result details: "+err.Error())
		return
	}
}

// SendErrorは，WebSocketを使用しているクライアントに対してエラーメッセージを送信し，その後コネクションを適切にクローズする関数である．
// この関数は，エラーが発生した際に，そのエラーメッセージをクライアントに通知するために使用される．
// エラーメッセージはテキストメッセージとして送信され，その送信が完了した後，WebSocketコネクションは正常終了のクローズメッセージを送ってクローズされる．
// エラーメッセージの送信やコネクションのクローズに失敗した場合，そのエラーはログに記録されるが，さらなるエラーハンドリングは行われない．
// この関数は，サーバーとクライアント間の通信中に発生する可能性のあるエラーをクライアントに伝達し，コネクションを清潔に保つために重要である．
//
// パラメータ:
// - conn *websocket.Conn: エラーメッセージを送信するWebSocketコネクション．
// - errMsg string: クライアントに送信するエラーメッセージの内容．
//
// 注意:
// - エラーメッセージの送信後，WebSocketコネクションは正常終了のためのクローズメッセージを送信してからクローズされる．
// - この関数の実行中に発生したエラーは，ログに記録されるが，呼び出し元には伝播されない．
func SendError(conn *websocket.Conn, errMsg string) {
	// エラーメッセージをクライアントに送信
	if err := conn.WriteMessage(websocket.TextMessage, []byte(errMsg)); err != nil {
		// エラーメッセージの送信に失敗した場合のログ出力やエラーハンドリング
		log.Printf("Failed to send error message: %v", err)
	}
	// WebSocketコネクションを適切にクローズ
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		// コネクションクローズの失敗に関するログ出力やエラーハンドリング
		log.Printf("Failed to close websocket connection: %v", err)
	}
}

func sendRequestToJudgeServer(ctx context.Context, url string, data []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// コンテキストキャンセルによるエラーを確認
		if errors.Is(err, context.Canceled) {

			return nil, fmt.Errorf("request to judge server was canceled due to timeout")
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
