package utils

import (
	"io"
	"mime/multipart"
	"net/http"
)

func IsImage(file multipart.File) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}
	file.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(buffer)

	return allowedTypes[contentType]
}
