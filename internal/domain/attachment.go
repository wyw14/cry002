package domain

import "time"

type Attachment struct {
	ID           string    `json:"id"`
	CaseID       string    `json:"case_id"`
	MaterialID   string    `json:"material_id,omitempty"`
	OriginalName string    `json:"original_name"`
	StorageKey   string    `json:"storage_key"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	SHA256       string    `json:"sha256"`
	Version      int       `json:"version"`
	UploadedBy   string    `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}
