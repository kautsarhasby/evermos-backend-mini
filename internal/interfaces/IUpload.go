package interfaces

import (
	"context"
	"errors"
	"io"
)

var ErrUploadTooLarge = errors.New("File exceeds 5 MB")

type IUploadService interface {
	Upload(ctx context.Context, filename string, file io.Reader) (string, error)
}
