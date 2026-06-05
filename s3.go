package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	client := s3.NewPresignClient(s3Client)

	req, err := client.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}
	bucket, key, found := strings.Cut((*video.VideoURL), ",")
	if found == false {
		return database.Video{}, errors.New("Invalid database URL")
	}

	expireTime := time.Hour * 72
	url, err := generatePresignedURL(cfg.s3Client, bucket, key, expireTime)
	if err != nil {
		return database.Video{}, err
	}

	video.VideoURL = &url
	return video, nil
}
