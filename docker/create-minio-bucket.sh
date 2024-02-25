#!/bin/sh

# MinIOクライアントの設定
mc config host add io_files http://minio:9000 "${MINIO_ROOT_USER}" "${MINIO_ROOT_PASSWORD}"

# バケットの作成
mc mb io_files/"${MINIO_BUCKET_NAME}"

# バケットのアクセスポリシーを公開設定(public)に設定
# minioバケットに関する外部からのアクセス権限の付与
mc anonymous set public io_files/"${MINIO_BUCKET_NAME}"
