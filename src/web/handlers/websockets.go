package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"procon_web_service/src/common/models"
	"procon_web_service/src/common/utils"
	"procon_web_service/src/web/async"
	webutils "procon_web_service/src/web/utils"
	"time"

	"github.com/gorilla/websocket"
)

// upgraderは，HTTPリクエストをWebSocketプロトコルにアップグレードするための設定を保持するwebsocket.Upgrader構造体のインスタンスである．
// このインスタンスは，WebSocket通信の際に使用され，ReadBufferSizeおよびWriteBufferSizeプロパティによって，読み取りと書き込みのバッファサイズが指定される．
// CheckOrigin関数は，WebSocket接続を試みるオリジンの検証を行う．ここでは，すべてのオリジンからの接続を許可するためにtrueを返している．
// これにより，異なるオリジンからのWebSocket接続が可能となり，クロスオリジンリソース共有(CORS)ポリシーによる制限を回避することができる．
//
// 使用例:
// WebSocket接続の要求を受け取ったサーバーは，このupgraderを使用してHTTPリクエストをWebSocket接続にアップグレードし，クライアントとの双方向通信を開始する．
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1 << 15,
	WriteBufferSize: 1 << 15,
	CheckOrigin: func(r *http.Request) bool {
		return true // すべてのオリジンを許可
	},
}

// WebSocketHandlerは，WebSocket通信を介して解答の提出と判定結果の非同期通知を処理するHTTPハンドラ関数である．
// クライアントからWebSocket接続が確立されると，この関数は接続をアップグレードし，クライアントとの間でメッセージを非同期にやり取りする準備をする．
// クライアントから送信された解答データを受け取り，非同期に判定処理を行う．解答データの判定は，JudgeSolutionAsync関数によって実行され，判定結果はWebSocketを通じてクライアントに通知される．
// WebSocket接続は，通信が完了するまで，またはエラーが発生するまで維持される．エラーが発生した場合，適切なエラーメッセージがクライアントに送信される．
// このハンドラは，Webサーバーとjudge-server間で別の通信メカニズム（例えばHTTPリクエスト）を使用して，解答の判定を非同期に行う設計になっている．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - func(w http.ResponseWriter, r *http.Request): WebSocket通信を処理する関数．
func WebSocketHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータからトークンを取得
		token := r.URL.Query().Get("token")

		// URLに渡されたトークンの検証
		if _, err := webutils.IsTokenAuthenticated(token); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err) // HTTPレスポンスの送信はここでは行わず，ログに記録するのみ
			return
		}
		defer conn.Close()

		// コンテキストの設定 : 時間制限: 2 minutes
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// クライアント接続エラー
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					async.SendError(conn, "WebSocket closed unexpectedly")
				}
				break
			}

			// 解答提出リクエストの処理
			var solution models.Solution
			if err := json.Unmarshal(message, &solution); err != nil {
				async.SendError(conn, "Failed to unmarshal solution: "+err.Error())
				continue
			}

			// 解答に基づいて非同期処理をトリガー
			go async.JudgeSolutionAsync(ctx, db, solution, conn)
		}
	}
}
