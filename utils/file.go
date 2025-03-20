package utils

import (
	"bytes"
	"context"
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

func GetFile(path string) ([]byte, error) {
	namespace := os.Getenv("BUCKET_NAMESPACE")
	bucket := os.Getenv("BUCKET_NAME")
	request := objectstorage.GetObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucket,
		ObjectName:    &path,
	}

	response, err := GetClient().GetObject(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("could not get object: %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from object on OCI : %v", err)
	}

	return buf.Bytes(), nil
}

func UploadFile(fileHeader *multipart.FileHeader, path string) (error, string) {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("could not open file: %v", err), ""
	}
	defer file.Close()

	extension := GetImageExtension(file)
	if extension == "" {
		return fmt.Errorf("invalid image extension"), ""
	}

	namespace := os.Getenv("BUCKET_NAMESPACE")
	bucket := os.Getenv("BUCKET_NAME")
	objectName := path + uuid.New().String() + extension

	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucket,
		ObjectName:    &objectName,
		PutObjectBody: file,
		ContentLength: &fileHeader.Size,
	}

	_, err = GetClient().PutObject(context.Background(), request)
	if err != nil {
		return err, ""
	}

	return nil, objectName
}
