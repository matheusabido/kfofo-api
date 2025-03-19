package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

func GetImageExtension(file multipart.File) string {
	allowedTypes := map[string]string{
		"image/jpeg": ".jpeg",
		"image/png":  ".png",
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return ""
	}
	file.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(buffer)

	return allowedTypes[contentType]
}

func CreateUploadRequest(fileHeader *multipart.FileHeader, path string) (*objectstorage.PutObjectRequest, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	extension := GetImageExtension(file)
	if extension == "" {
		return nil, fmt.Errorf("invalid image extension")
	}

	namespace := os.Getenv("BUCKET_NAMESPACE")
	bucket := os.Getenv("BUCKET_NAME")
	objectName := path + "/" + uuid.New().String() + extension

	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucket,
		ObjectName:    &objectName,
		PutObjectBody: file,
		ContentLength: &fileHeader.Size,
	}

	return &request, nil
}
