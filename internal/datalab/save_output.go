package datalab

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/minio/minio-go"
)

func SaveRecords(endpoint string, accessKeyID string, secretAccessKey string, rec []types.Log, bucket, key string) error {
	useSSL := true
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", minioClient)

	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	reader := strings.NewReader(string(data))

	info, err := minioClient.PutObject(bucket, key, reader, reader.Size(), minio.PutObjectOptions{ContentType: "text/plain"})

	if err != nil {
		return err
	}
	fmt.Printf("%v\n", info)

	return nil
}
