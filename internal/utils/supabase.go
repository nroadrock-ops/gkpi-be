package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"gkpi-be/internal/config"
)

// UploadToSupabaseStorage mengunggah file ke bucket Supabase dan mengembalikan public URL-nya
func UploadToSupabaseStorage(file *multipart.FileHeader, bucketName string) (string, error) {
	cfg := config.LoadConfig()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	// Ganti spasi dengan underscore untuk keamanan URL
	safeFileName := strings.ReplaceAll(file.Filename, " ", "_")
	
	// URL API Supabase untuk upload (Storage API)
	uploadPath := fmt.Sprintf("%s/storage/v1/object/%s/%s", cfg.SupabaseURL, bucketName, safeFileName)

	req, err := http.NewRequest("POST", uploadPath, bytes.NewBuffer(fileBytes))
	if err != nil {
		return "", err
	}

	// Gunakan Service Key untuk otorisasi bypass RLS
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceKey)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Coba baca body error dari Supabase untuk pesan yang lebih jelas
		bodyErr, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload file, status: %d, msg: %s", resp.StatusCode, string(bodyErr))
	}

	// Generate public URL yang benar menggunakan env SUPABASE_URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", cfg.SupabaseURL, bucketName, safeFileName)

	return publicURL, nil
}
