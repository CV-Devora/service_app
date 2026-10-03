package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	v1 "toko-emas/api/v1"

	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file too large or invalid form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file is required"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file type not allowed (jpg, png, gif, webp only)"})
		return
	}

	filename := uuid.New().String() + ext
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, 0755)
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: "failed to write file"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", filename)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "upload success", Data: fileURL})
}
