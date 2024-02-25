package models

import "time"

// ResultDetailは，解答の評価結果の詳細を表す構造体である．
// これには，テストケースの総数，正解数，不正解数，実行時間超過数，各テストケースの結果，
// およびエラーメッセージが含まれる．
type ResultDetail struct {
	TotalCases        int          `json:"total_cases"`         // テストケースの総数である．
	CorrectCases      int          `json:"correct_cases"`       // 正解したテストケースの数である．
	IncorrectCases    int          `json:"incorrect_cases"`     // 不正解だったテストケースの数である．
	TimeLimitExceeded int          `json:"time_limit_exceeded"` // 実行時間超過となったテストケースの数である．
	CaseResults       []CaseResult `json:"case_results"`        // 各テストケースの詳細結果を含む配列である．
	ErrorMessage      string       `json:"error,omitempty"`     // 解答の実行中に発生したエラーメッセージである（存在する場合）．
}

// CaseResultは，個々のテストケースの実行結果を表す構造体である．
// テストケースの名前，結果（正解，不正解，実行時間超過），および実行時間が含まれる．
type CaseResult struct {
	CaseName      string        `json:"case_name"`      // テストケースの名前である．
	Result        string        `json:"result"`         // テストケースの結果（"correct", "incorrect", "time_limit_exceeded"）である．
	ExecutionTime time.Duration `json:"execution_time"` // テストケースの実行時間（ミリ秒）である．
}

// Solutionは，ユーザーが提出した解答の情報を保持する構造体である．
type Solution struct {
	SolutionID  int       `json:"solution_id"`  // 解答の一意識別子である．
	UserID      int       `json:"user_id"`      // 解答を提出したユーザーのIDである．
	ProblemID   int       `json:"problem_id"`   // 解答が対象とする問題のIDである．
	LanguageID  int       `json:"language_id"`  // 解答が記述されたプログラミング言語のIDである．
	Code        string    `json:"code"`         // 解答のソースコードである．
	SubmittedAt time.Time `json:"submitted_at"` // 解答の提出日時である．
}
