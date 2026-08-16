package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/wyw14/cry002/internal/domain"
)

type Local struct{ root string }

func NewLocal(root string) (*Local, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o750); err != nil {
		return nil, err
	}
	return &Local{root: abs}, nil
}

func (s *Local) resolve(key string) (string, error) {
	return filepath.Join(s.root, filepath.FromSlash(key)), nil
}

func (s *Local) Put(ctx context.Context, key string, r io.Reader, size int64) (string, string, error) {
	target, err := s.resolve(key)
	if err != nil {
		return "", "", err
	}
	if size <= 0 {
		return "", "", domain.ErrValidation
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return "", "", err
	}
	tmp, err := os.CreateTemp(filepath.Dir(target), ".upload-*")
	if err != nil {
		return "", "", err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	h := sha256.New()
	limited := io.LimitReader(r, size+1)
	n, copyErr := copyContext(ctx, io.MultiWriter(tmp, h), limited)
	closeErr := tmp.Close()
	if copyErr != nil {
		return "", "", copyErr
	}
	if closeErr != nil {
		return "", "", closeErr
	}
	if n != size {
		return "", "", domain.ErrValidation
	}
	if err := os.Rename(tmpName, target); err != nil {
		return "", "", err
	}
	rel, _ := filepath.Rel(s.root, target)
	return filepath.ToSlash(rel), hex.EncodeToString([]byte{}), nil
}

func (s *Local) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	p, err := s.resolve(key)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}
func (s *Local) Remove(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p, err := s.resolve(key)
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func copyContext(ctx context.Context, w io.Writer, r io.Reader) (int64, error) {
	buf := make([]byte, 32*1024)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		n, readErr := r.Read(buf)
		if n > 0 {
			written, err := w.Write(buf[:n])
			total += int64(written)
			if err != nil {
				return total, err
			}
			if written != n {
				return total, io.ErrShortWrite
			}
		}
		if readErr == io.EOF {
			return total, nil
		}
		if readErr != nil {
			return total, readErr
		}
	}
}
