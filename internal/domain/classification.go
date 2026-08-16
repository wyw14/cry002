package domain

import "time"

type Classification struct {
	ID        string    `json:"id"`
	ParentID  string    `json:"parent_id,omitempty"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	SortOrder int       `json:"sort_order"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateClassificationMove(nodeID, parentID string, nodes map[string]Classification) error {
	if nodeID == "" || nodeID == parentID {
		return ErrValidation
	}
	seen := map[string]bool{nodeID: true}
	for cur := parentID; cur != ""; {
		if seen[cur] {
			return ErrValidation
		}
		seen[cur] = true
		n, ok := nodes[cur]
		if !ok {
			return ErrNotFound
		}
		cur = n.ParentID
	}
	return nil
}
