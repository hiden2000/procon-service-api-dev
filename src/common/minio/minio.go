package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"procon_web_service/src/common/config"
	commonerrors "procon_web_service/src/common/errors"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	bucketName  = ""
	minIOClient *minio.Client
	minioconfig = config.NewMinIOConfig()
)

// initはMinIOクライアントを初期化する．
// MinIOの設定は環境変数から読み込まれ，クライアントが正常に初期化されるとグローバル変数に設定される．
// 初期化に失敗した場合はアプリケーションを終了させる．
func init() {
	var err error
	minIOClient, err = initializeMinIOClient()
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}
	bucketName = minioconfig.BucketName
}

// initializeMinIOClientはMinIOクライアントを初期化するユーティリティ関数である．
// MinIOのエンドポイント，認証情報，およびSSLの使用有無を設定してMinIOクライアントを生成する．
// 成功した場合は生成されたMinIOクライアントを返し，失敗した場合はエラーを返す．
func initializeMinIOClient() (*minio.Client, error) {
	minioClient, err := minio.New(minioconfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioconfig.User, minioconfig.Password, ""),
		Secure: minioconfig.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

// UploadFileToMinIOは，指定された問題IDに関連するファイルをMinIOにアップロードする．
//
// この関数は，提供された複数のmultipart.FileHeader（fileHeaders）からファイルを読み込み，
// 指定されたfileType（'in'や'out'など）に基づいてMinIO内の適切な場所にアップロードする．
// アップロード中にエラーが発生した場合，既にアップロードされたファイルはクリーンアップされ，エラーが返される．
// 全てのファイルのアップロードが成功した場合は，nilエラーが返される．
//
// パラメータ:
// - problemID int: アップロードされるファイルが関連する問題のID．
// - fileHeaders []*multipart.FileHeader: アップロードするファイルのヘッダのスライス．
// - fileType string: アップロードされるファイルのタイプ（例：'in'，'out'）．
//
// 戻り値:
// - error: アップロード中に発生したエラー，またはnil．
func UploadFileToMinIO(problemID int, fileHeaders []*multipart.FileHeader, fileType string) error {
	uploadedFilePaths := []string{} // アップロードされたファイルのパスを追跡

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			// エラー時には既にアップロードされたファイルを削除
			cleanupUploadedFiles(uploadedFilePaths)
			return commonerrors.WrapMinIOError("uploading files to MinIO", err)
		}
		defer file.Close()

		filePath := GetFileSaveName("", problemID, fileType, fileHeader.Filename)

		// ファイルのアップロード
		_, err = minIOClient.PutObject(context.Background(), bucketName, filePath, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType:        "text/plain",
			ContentEncoding:    "utf-8",
			ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", fileHeader.Filename),
		})
		if err != nil {
			// エラー時には既にアップロードされたファイルを削除
			cleanupUploadedFiles(uploadedFilePaths)
			return commonerrors.WrapMinIOError("uploading files to MinIO", err)
		}

		uploadedFilePaths = append(uploadedFilePaths, filePath)
	}

	return nil
}

// DownloadIOFilesは，指定された問題IDに関連する入力ファイルと出力ファイルをMinIOからダウンロードし，ローカルの/tmpディレクトリに保存する．
//
// この関数は，MinIOから'in'および'out'ファイルタイプに関連するファイルをダウンロードし，
// 指定された問題IDに基づいてローカルの保存先ディレクトリを生成する．ダウンロードまたはディレクトリの作成中にエラーが発生した場合，エラーが返される．
//
// パラメータ:
// - ctx context.Context: 操作のコンテキスト．
// - problemID int: ダウンロードするファイルが関連する問題のID．
//
// 戻り値:
// - error: ダウンロードまたはディレクトリの作成中に発生したエラー，またはnil．
func DownloadIOFiles(ctx context.Context, problemID int) error {
	savedDirName := GetFileSaveName("/tmp", problemID, "", "")

	if err := os.MkdirAll(savedDirName, 0755); err != nil {
		return commonerrors.WrapMinIOError("", err)
	}

	// MinIOから入力ファイルと出力ファイルをダウンロード
	if err := downloadFilesFromMinIO(ctx, problemID, "in", savedDirName); err != nil {
		return commonerrors.WrapMinIOError("downloading files from MinIO", err)
	}
	if err := downloadFilesFromMinIO(ctx, problemID, "out", savedDirName); err != nil {
		return commonerrors.WrapMinIOError("downloading files from MinIO", err)
	}

	return nil
}

// DeleteFileFromMinIOは，MinIOの特定のバケットから，指定されたプレフィックスを持つ全てのファイルを削除する．
//
// この関数は，MinIO内のbucketNameバケットから，指定されたプレフィックス（prefix）に一致する
// 全てのファイルを検索し，それらを削除する．ファイルの検索または削除中にエラーが発生した場合，エラーが返される．
//
// パラメータ:
// - prefix string: 削除するファイルのプレフィックス．
//
// 戻り値:
// - error: ファイルの削除中に発生したエラー，またはnil．
func DeleteFileFromMinIO(prefix string) error {

	// バケット内の指定されたプレフィックスを持つオブジェクトのリストを取得
	objectCh := minIOClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	// リストからオブジェクトを取得し削除
	for object := range objectCh {
		if object.Err != nil {
			return commonerrors.WrapMinIOError("cleaning files to MinIO", object.Err)
		}
		err := minIOClient.RemoveObject(context.Background(), bucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return commonerrors.WrapMinIOError("cleaning files to MinIO", err)
		}
	}

	return nil
}

// GetFileSaveNameはファイルを保存する際のフォルダ名を生成する関数である．
// 指定されたディレクトリ名，問題ID，ファイルタイプ，ファイル名から成るパスを返す．
// この関数はファイルの整理とアクセスのための一貫した命名規則を提供する．
func GetFileSaveName(dirName string, problemID int, fileType, fileName string) string {
	return filepath.Join(dirName, "problem_"+strconv.Itoa(problemID), fileType, fileName)
}

// downloadFilesFromMinIOは特定の問題IDとファイルタイプに基づいてMinIOからファイルをダウンロードし，
// 指定されたローカルディレクトリに保存する内部関数である．
// コンテキスト，問題ID，ファイルタイプ，および保存先ディレクトリ名を受け取る．
// ファイルのダウンロードまたは保存中にエラーが発生した場合は，エラーを返す．
func downloadFilesFromMinIO(ctx context.Context, problemID int, fileType, savedDirName string) error {
	prefix := GetFileSaveName("", problemID, fileType, "")

	if err := os.MkdirAll(filepath.Join(savedDirName, fileType), 0755); err != nil {
		return fmt.Errorf("failed to create directory for IO files: %w", err)
	}

	for object := range minIOClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return fmt.Errorf("error listing object: %w", object.Err)
		}

		savePath := filepath.Join(savedDirName, fileType, filepath.Base(object.Key))

		// ファイルが既に存在するか + 最新であるか確認し必要に応じてダウンロード
		if needDownload(savePath, object.LastModified) {
			if err := downloadAndSaveFile(ctx, object.Key, savePath); err != nil {
				return fmt.Errorf("failed to download and save file: %w", err)
			}
		}
	}

	return nil
}

// needDownloadは指定されたファイルがダウンロードが必要かどうかを判断する関数である．
// ファイルパスとMinIOのファイルの最終更新時刻を受け取り，
// ローカルファイルが存在しない，またはMinIOのファイルよりも古い場合はtrueを返す．
func needDownload(filePath string, lastModified time.Time) bool {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// ファイルが存在しない場合はダウンロードが必要
		return true
	} else if err != nil {
		// ファイルの状態取得中にエラーが発生した場合安全のためにダウンロードが必要と判断
		log.Printf("Error checking file info: %v", err)
		return true
	}

	// ローカルファイルの最終更新時刻がMinIOのファイルの最終更新時刻よりも古いかどうかをチェック
	return fileInfo.ModTime().Before(lastModified)
}

// downloadAndSaveFileはMinIOから特定のファイルをダウンロードし，指定されたパスに保存する関数である．
// コンテキスト，MinIO内のオブジェクトキー，および保存先のファイルパスを受け取る．
// ダウンロードまたは保存中にエラーが発生した場合は，エラーを返す．
func downloadAndSaveFile(ctx context.Context, objectKey, savePath string) error {
	// MinIOからファイルをダウンロードしてsavePathに保存
	object, err := minIOClient.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object from MinIO: %v, object key: %s", err, objectKey)
	}
	defer object.Close()

	// ダウンロードしたファイルをローカルに保存
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file for saving object: %v, path: %s", err, savePath)
	}
	defer file.Close()

	if _, err = io.Copy(file, object); err != nil {
		return fmt.Errorf("failed to write object to file: %v, path: %s", err, savePath)
	}

	return nil
}

// cleanupUploadedFilesはアップロードされたが不要になったファイルをMinIOから削除する関数である．
// アップロードされたファイルのパスのスライスを受け取り，それらのファイルをMinIOから削除する．
// ファイルの削除中にエラーが発生した場合はログに記録するが，プロセスを中断しない．
func cleanupUploadedFiles(filePaths []string) {
	for _, filePath := range filePaths {
		err := minIOClient.RemoveObject(context.Background(), bucketName, filePath, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Failed to delete uploaded file: %s, error: %v", filePath, err)
		}
	}
}
