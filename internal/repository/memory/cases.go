package memory

import (
	"context"
	"sort"
	"strings"

	"github.com/wyw14/cry002/internal/application"
	"github.com/wyw14/cry002/internal/domain"
)

func (s *Store) CreateCase(ctx context.Context, c domain.CaseFile, materials []domain.Material) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.cases {
		if existing.CaseNumber == c.CaseNumber {
			return domain.ErrConflict
		}
	}
	s.cases[c.ID] = c
	s.materials[c.ID] = append([]domain.Material(nil), materials...)
	return nil
}

func (s *Store) UpdateCase(ctx context.Context, c domain.CaseFile, materials []domain.Material) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cases[c.ID]; !ok {
		return domain.ErrNotFound
	}
	s.cases[c.ID] = c
	s.materials[c.ID] = append([]domain.Material(nil), materials...)
	return nil
}

func (s *Store) CaseByID(ctx context.Context, id string) (domain.CaseFile, error) {
	if err := check(ctx); err != nil {
		return domain.CaseFile{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.cases[id]
	if !ok {
		return domain.CaseFile{}, domain.ErrNotFound
	}
	return c, nil
}

func (s *Store) ListCases(ctx context.Context, f application.CaseFilter) ([]domain.CaseFile, int, error) {
	if err := s.wait(ctx); err != nil {
		return nil, 0, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]domain.CaseFile, 0, len(s.cases))
	q := strings.ToLower(strings.TrimSpace(f.Query))
	for _, c := range s.cases {
		if f.Status != "" && c.Status != f.Status {
			continue
		}
		if f.ClassificationID != "" && c.ClassificationID != f.ClassificationID {
			continue
		}
		if f.Department != "" && c.ResponsibleDepartment != f.Department {
			continue
		}
		if f.Year != 0 && c.Year != f.Year {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(c.CaseNumber+" "+c.Title+" "+c.Summary+" "+strings.Join(c.Keywords, " ")), q) {
			continue
		}
		items = append(items, c)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	total := len(items)
	page, size := f.Page, f.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= total {
		return []domain.CaseFile{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return append([]domain.CaseFile(nil), items[start:end]...), total, nil
}

func (s *Store) Materials(ctx context.Context, caseID string) ([]domain.Material, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Material(nil), s.materials[caseID]...), nil
}
func (s *Store) CreateVersion(ctx context.Context, v domain.CaseVersion) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.versions[v.CaseID] = append(s.versions[v.CaseID], v)
	return nil
}
func (s *Store) Versions(ctx context.Context, id string) ([]domain.CaseVersion, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.CaseVersion(nil), s.versions[id]...), nil
}
func (s *Store) CreateReview(ctx context.Context, v domain.Review) error {
	if err := check(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reviews[v.CaseID] = append(s.reviews[v.CaseID], v)
	return nil
}
func (s *Store) Reviews(ctx context.Context, id string) ([]domain.Review, error) {
	if err := check(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Review(nil), s.reviews[id]...), nil
}
