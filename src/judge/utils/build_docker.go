package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"procon_web_service/src/common/config"
	"procon_web_service/src/common/minio"
	"procon_web_service/src/common/models"
	"regexp"
	"strconv"
	"sync"
	"time"

	"strings"
)

var (
	// リソース制限の設定
	memoryLimit = "512m" // メモリ制限
	cpuLimit    = "1.0"  // CPU制限
	timeout     = 1      //実行時間制限

	// コンテナ内部のセキュリティオプション指定
	securityOpts = []string{
		"--net=none",                      // コンテナのネットワークアクセスを遮断
		"--read-only",                     // コンテナのファイルシステムを読み取り専用に設定
		"--tmpfs /workspace:rw,size=512m", // 読み書き可能で512MBのサイズ制限付きの一時ファイルシステムとしてマウント
		"--cap-drop=ALL",                  // セキュリティのためコンテナ内部Linuxカーネル機能を削除
	}
)

type ExecutionResult struct {
	Success       bool
	ExecutionTime int64  // 実行時間（ナノ秒）
	OutputDiff    bool   // 出力が期待される出力と異なる場合はtrue
	ErrorMessage  string // 実行エラーのメッセージ（エラーが発生した場合）
}

// BuildAndRunInContainer - 提出されたコードをDockerコンテナ内で平行処理によりテスト && 結果を取得
func BuildAndRunInContainer(ctx context.Context, solution models.Solution) (*models.ResultDetail, error) {
	_, ok := config.GetLanguageConfigByID(solution.LanguageID)
	if !ok {
		return nil, fmt.Errorf("unsupported language ID: %d", solution.LanguageID)
	}

	// ソースコードを一時ファイルに保存
	codeFilePath, langConfig, cleanup, err := SaveCodeToFile(solution.LanguageID, solution.Code)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	// 関連する入出力ファイルをダウンロードしローカルに保存
	if err := minio.DownloadIOFiles(ctx, solution.ProblemID); err != nil {
		return nil, err
	}

	ioFiles, err := getTestCases(solution.ProblemID)
	if err != nil {
		return nil, err
	}

	var results models.ResultDetail
	results.TotalCases = len(ioFiles)

	errChan := make(chan error, 1)                            // エラーを受け取るチャネル
	resultsChan := make(chan models.CaseResult, len(ioFiles)) // 結果を受け取るチャネル

	var wg sync.WaitGroup
	// 並列処理によってテストケースを実行
	// ATTENTION: 同時に多数のDockerコンテナを起動することになるためシステムリソースに与える影響を考慮し適切な並列度を設定する必要あり
	for inputFilePath, outputFilePath := range ioFiles {

		wg.Add(1)
		go func(inputFilePath, outputFilePath string) {
			defer wg.Done()

			// 一時的な書き込みファイルを作成
			tempFilePath, cleanup, err := CreateTempFile()
			if err != nil {
				select {
				case errChan <- err:
				case <-ctx.Done():
				}
				return
			}
			defer cleanup()

			dockerCommand := buildDockerRunCommandWithTimeout(langConfig, codeFilePath, inputFilePath, tempFilePath, outputFilePath, solution.ProblemID, timeout)
			executionResult, err := executeDockerCommand(ctx, dockerCommand)
			if err != nil {
				select {
				case errChan <- err:
				case <-ctx.Done():
				}
				return
			}

			caseResult := models.CaseResult{
				CaseName: filepath.Base(inputFilePath),
				Result: func() string {
					if executionResult.ErrorMessage != "" {
						results.IncorrectCases++
						return "INTERNAL ERROR"
					}
					if executionResult.OutputDiff {
						results.IncorrectCases++
						return "FAILED"
					}
					if executionResult.ExecutionTime/1_000_000 > 2000 {
						results.TimeLimitExceeded++
						return "TIME LIMITED EXCEEDED"
					}
					results.CorrectCases++
					return "PASSED"
				}(),
				ExecutionTime: time.Duration(executionResult.ExecutionTime),
			}

			select {
			case resultsChan <- caseResult:
			case <-ctx.Done():
			}
		}(inputFilePath, outputFilePath)

	}

	// ゴルーチンがすべて終了するのを待つ
	go func() {
		wg.Wait()
		close(errChan)
		close(resultsChan)
	}()

	// 結果とエラーの処理
	for {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		case result, ok := <-resultsChan:
			if !ok {
				// チャネルが閉じられたとき，全てのケースが処理されたと判断
				return &results, nil
			}
			results.CaseResults = append(results.CaseResults, result)
		case <-ctx.Done():
			// コンテキストがキャンセルされた場合，処理を中断
			return nil, ctx.Err()
		}
	}
}

// buildDockerRunCommandWithTimeout - 提出された言語設定に基づいて適切なDockerコマンドを構築
func buildDockerRunCommandWithTimeout(langConfig config.LanguageConfig, codeFilePath, inputFilePath, tempFilePath, outputFilePath string, problemID int, timeout int) string {
	// コードファイルのマウント設定
	codeFileVolume := fmt.Sprintf("-v %s:/workspace/code", filepath.Dir(codeFilePath))

	// 入出力ファイルのマウント設定
	ioVolumeMapping := fmt.Sprintf("-v %s:/workspace/io", minio.GetFileSaveName("/tmp", problemID, "", ""))

	// 一時ファイルの保存先を /workspace/tmp に設定
	tempFileVolume := fmt.Sprintf("-v %s:/workspace/tmp", filepath.Dir(tempFilePath))

	// コンパイルコマンドの準備
	compileCmd := ""
	if langConfig.Compile != "" {
		compileCmd = strings.Replace(langConfig.Compile, "{code}", "/workspace/code/"+filepath.Base(codeFilePath), -1) + " && "
	}

	// 実行コマンドの準備
	runCmd := strings.Replace(langConfig.Run, "{code}", "/workspace/code/"+filepath.Base(codeFilePath), -1)

	// 入力ファイルから一時ファイルへのリダイレクトを含む実行コマンド
	executionCmd := fmt.Sprintf("%s < /workspace/io/in/%s > /workspace/tmp/%s", runCmd, filepath.Base(inputFilePath), filepath.Base(tempFilePath))

	// タイムアウトと実行時間計測を含む実行コマンド(コンパイル時の一時ファイルの作成場所を変更)
	commandWithTimeout := fmt.Sprintf("%s %s timeout --preserve-status %ds /bin/sh -c 'start=\\$(date +%%s%%N); %s; end=\\$(date +%%s%%N); echo Execution time: \\$((end-start)) nanoseconds'", langConfig.Setup, compileCmd, timeout, executionCmd)

	// 結果比較スクリプトの追加（一時ファイルと期待される出力ファイルを比較）
	compareScript := fmt.Sprintf(" && diff -q /workspace/tmp/%s /workspace/io/out/%s", filepath.Base(tempFilePath), filepath.Base(outputFilePath))

	// Dockerコマンドの組み立て
	dockerCommand := fmt.Sprintf("docker run --rm %s %s %s %s --memory %s --cpus %s %s /bin/sh -c \"%s%s\"",
		strings.Join(securityOpts, " "), codeFileVolume, ioVolumeMapping, tempFileVolume, memoryLimit, cpuLimit, langConfig.Image, commandWithTimeout, compareScript)

	return dockerCommand
}

// getTestCases - 指定された問題IDに基づき入出力ファイルのパスのペアを返す
func getTestCases(problemID int) (map[string]string, error) {
	testCases := make(map[string]string)
	inDir := minio.GetFileSaveName("/tmp", problemID, "in", "")
	outDir := minio.GetFileSaveName("/tmp", problemID, "out", "")

	// 入力ファイルのリストを取得
	inputFiles, err := ioutil.ReadDir(inDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read input directory: %v", err)
	}

	// 出力ファイルのリストを取得
	if _, err := ioutil.ReadDir(outDir); err != nil {
		return nil, fmt.Errorf("failed to read output directory: %v", err)
	}

	// 入力ファイルと出力ファイルをマッピング
	for _, inFile := range inputFiles {
		if !inFile.IsDir() { // ディレクトリではないことを確認
			inFilePath := filepath.Join(inDir, inFile.Name())
			outFilePath := filepath.Join(outDir, inFile.Name())
			// 出力ファイルが存在するか確認
			if _, err := ioutil.ReadFile(outFilePath); err == nil {
				testCases[inFilePath] = outFilePath
			} else {
				return nil, fmt.Errorf("output file does not exist for input file %s", inFile.Name())
			}
		}
	}

	return testCases, nil
}

// executeDockerCommand - Dockerコマンドを実行 && 結果を解析
func executeDockerCommand(ctx context.Context, dockerCmd string) (ExecutionResult, error) {
	// コマンドを実行
	cmd := exec.CommandContext(ctx, "sh", "-c", dockerCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	result := ExecutionResult{}

	// コマンドの出力を解析
	output := out.String()

	// 実行時間を抽出
	timeRegex := regexp.MustCompile(`Execution time: (\d+) nanoseconds`)
	matches := timeRegex.FindStringSubmatch(output)
	if len(matches) > 1 {
		// matches[1]に実行時間が含まれる
		result.ExecutionTime, _ = strconv.ParseInt(matches[1], 10, 64)
	}

	// diffコマンドの結果を確認（出力が異なる場合はdiffコマンドから何らかの出力がある）
	diffRegex := regexp.MustCompile(`Files /workspace/tmp/.+ and /workspace/io/out/.+ differ`)
	if diffRegex.MatchString(output) {
		result.OutputDiff = true
	} else {
		result.OutputDiff = false
	}

	// エラーメッセージの設定
	if err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
	} else {
		result.Success = true
	}

	return result, nil
}
