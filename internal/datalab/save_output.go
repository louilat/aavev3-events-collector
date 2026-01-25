package datalab

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/minio/minio-go"
)

func SaveRecordsMINIO(endpoint string, accessKeyID string, secretAccessKey string, rec []types.Log, bucket, key string) error {
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

// SaveRecords uploads Ethereum logs to AWS S3 as JSON.
func SaveRecordsAWS(region, accessKeyID, secretAccessKey, bucket, key string, rec []types.Log) error {
	// 1. Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     accessKeyID,
					SecretAccessKey: secretAccessKey,
				}, nil
			},
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 2. Create S3 client
	client := s3.NewFromConfig(cfg)

	// 3. Marshal records to JSON
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("failed to marshal records: %w", err)
	}

	// 4. Upload to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to S3: %w", err)
	}

	fmt.Printf("Uploaded %d logs to s3://%s/%s\n", len(rec), bucket, key)
	return nil
}

func SaveRecordsLOCAL(path string, rec []types.Log) error {
	// 1. Ensure parent directories exist
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 2. Create (or truncate) the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. Encode data as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // optional: pretty JSON

	if err = encoder.Encode(rec); err != nil {
		return err
	}

	return nil
}
