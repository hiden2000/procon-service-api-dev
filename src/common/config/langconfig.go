package config

// LanguageConfigは言語ごとの実行環境設定を保持する構造体である．
// 各フィールドは特定のプログラミング言語におけるDockerイメージや実行コマンドなどを定義する．
type LanguageConfig struct {
	LanguageID int    // 言語の一意識別子．
	Image      string // 使用するDockerイメージの名前．
	Compile    string // ソースコードをコンパイルするためのコマンド．コンパイルが不要な場合は空文字列．
	Setup      string // 実行環境のセットアップに使用するコマンド．必要な環境変数の設定などを含む．
	Run        string // コンパイル済みのプログラム，またはスクリプトを実行するためのコマンド．
}

// SupportedLanguagesはサポートされる各言語の設定をLanguageIDをキーとして保持するマップである．
// このマップを通じて，言語ごとの実行環境設定にアクセスできる．
var SupportedLanguages = map[int]LanguageConfig{
	1: { // Python
		LanguageID: 1,
		Image:      "python:3.8-slim",
		Compile:    "",
		Run:        "python3 {code}",
	},
	2: { // C++
		LanguageID: 2,
		Image:      "gcc:latest",
		Compile:    "g++ {code} -o /workspace/tmp/a.out",
		Setup:      "export TMPDIR=/workspace/tmp &&",
		Run:        "/workspace/tmp/a.out",
	},
	3: { // Go
		LanguageID: 3,
		Image:      "golang:latest",
		Compile:    "go build -o /workspace/tmp/a.out {code}",
		Setup:      "export TMPDIR=/workspace/tmp && export GOCACHE=/workspace/tmp/go-cache &&",
		Run:        "/workspace/tmp/a.out",
	},
	4: { // Java
		LanguageID: 4,
		Image:      "openjdk:11",
		Compile:    "javac {code}",
		Setup:      "export TMPDIR=/workspace/tmp &&",
		Run:        "java -cp /workspace/code Main",
	},
	5: { // Rust
		LanguageID: 5,
		Image:      "rust:latest",
		Compile:    "rustc {code} -o /workspace/tmp/code",
		Setup:      "export TMPDIR=/workspace/tmp &&",
		Run:        "/workspace/tmp/code",
	},
}

// GetLanguageConfigByIDは指定されたLanguageIDに対応するLanguageConfigを返す関数である．
// 指定されたIDの設定が存在する場合はその設定とtrueを，存在しない場合はfalseを返す．
func GetLanguageConfigByID(languageID int) (LanguageConfig, bool) {
	config, ok := SupportedLanguages[languageID]
	return config, ok
}
