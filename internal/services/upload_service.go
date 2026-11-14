package services

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"evermos-hasby/ecommerce/internal/interfaces"
	"github.com/google/uuid"
)

const maxUploadSize = 5 << 20

type UploadService struct{}

func (s *UploadService) Upload(ctx context.Context, filename string, file io.Reader) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	name := uuid.NewString() + filepath.Ext(filepath.Base(filename))
	path := filepath.Join("uploads", name)
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", err
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.Remove(path)
		}
	}()

	n, err := io.Copy(dst, io.LimitReader(file, maxUploadSize+1))
	closeErr := dst.Close()
	if err != nil {
		return "", err
	}
	if n > maxUploadSize {
		return "", interfaces.ErrUploadTooLarge
	}
	if closeErr != nil {
		return "", closeErr
	}
	keep = true
	return name, nil
}
