package helper

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// Helper function to decode base64 and save image
func SaveBase64Image(ImageDir, base64String string) (string, error) {
	data := strings.Split(base64String, ",")
	if len(data) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}

	decoded, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return "", err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	filePath := filepath.Join(ImageDir, filename)

	err = ioutil.WriteFile(filePath, decoded, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
