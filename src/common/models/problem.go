package models

import "time"

// Problemは，コーディング問題の情報を保持する構造体である．
type Problem struct {
	ProblemID   int       `json:"problem_id"`   // 問題の一意識別子である．
	UserID      int       `json:"user_id"`      // 問題を作成したユーザーのIDである．
	Title       string    `json:"title"`        // 問題のタイトルである．
	Description string    `json:"description"`  // 問題の説明文である．
	Difficulty  int       `json:"difficulty"`   // 問題の難易度を表す整数値である．
	CreatedAt   time.Time `json:"created_at"`   // 問題の作成日時である．
	UpdatedAt   time.Time `json:"updated_at"`   // 問題の最終更新日時である．
	CategoryIDs []int     `json:"category_ids"` // 問題に関連付けられたカテゴリIDのリストである．
}
