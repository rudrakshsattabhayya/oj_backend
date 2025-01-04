package helpers

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
)

func UploadFile(file *multipart.FileHeader, bucketName string) (string, error) {
	supabaseClient := config.GetSupabaseClient()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	filePath := uuid.New().String() + "/" + file.Filename

	_, err = supabaseClient.UploadFile(bucketName, filePath, src)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %s", err.Error())
	}

	publicURL := supabaseClient.GetPublicUrl(bucketName, filePath).SignedURL

	return publicURL, nil
}
