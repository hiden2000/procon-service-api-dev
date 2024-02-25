package handlers

import (
	"errors"
	"net/http"
	"procon_web_service/src/common/models"
	"procon_web_service/src/common/utils"
	judgeutils "procon_web_service/src/judge/utils"
)

// JudgeHandler - リクエストに添付された伝播contextを利用して非同期通信をコントロール
func JudgeHandler(w http.ResponseWriter, r *http.Request) {
	var solution models.Solution

	if err := utils.DecodeRequestBody(r, &solution); err != nil {
		utils.SendErrorResponse(w, err)
		return
	}

	// http.Requestのコンテキストを取得
	ctx := r.Context()

	resultChan := make(chan *models.ResultDetail)
	errChan := make(chan error)
	go func() {
		resultDetail, err := judgeutils.BuildAndRunInContainer(ctx, solution)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- resultDetail
	}()

	select {
	case resultDetail := <-resultChan:
		if resultDetail.ErrorMessage != "" {
			utils.SendErrorResponse(w, errors.New(resultDetail.ErrorMessage))
			return
		}
		utils.SendJSONResponse(w, http.StatusOK, *resultDetail)
	case err := <-errChan:
		utils.SendErrorResponse(w, err)
	case <-ctx.Done():
		utils.SendErrorResponse(w, errors.New("request timed out"))
	}
}
