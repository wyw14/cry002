package tests

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/wyw14/cry002/internal/platform/storage"
	"io"
	"testing"
)

func TestAttachmentStorageRejectsTraversalAndPreservesHash(t *testing.T) {
	s, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.Put(context.Background(), "../escape.pdf", bytes.NewReader([]byte("x")), 1); err == nil {
		t.Fatal("path traversal accepted")
	}
	body := []byte("verified archive attachment")
	key, sum, err := s.Put(context.Background(), "case/version/file.pdf", bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	want := sha256.Sum256(body)
	if sum != hex.EncodeToString(want[:]) {
		t.Fatalf("hash=%s", sum)
	}
	r, err := s.Open(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	got, _ := io.ReadAll(r)
	if !bytes.Equal(got, body) {
		t.Fatal("stored content changed")
	}
}
