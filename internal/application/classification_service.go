package application

import (
	"context"
	"github.com/wyw14/cry002/internal/domain"
	"sort"
)

type ClassificationService struct {
	repo  Repository
	clock Clock
	ids   IDGenerator
}

func NewClassificationService(repo Repository, clock Clock, ids IDGenerator) *ClassificationService {
	return &ClassificationService{repo: repo, clock: clock, ids: ids}
}
func (s *ClassificationService) Tree(ctx context.Context) ([]domain.Classification, error) {
	items, err := s.repo.ListClassifications(ctx)
	sort.SliceStable(items, func(i, j int) bool { return items[i].SortOrder < items[j].SortOrder })
	return items, err
}
func (s *ClassificationService) Create(ctx context.Context, actor domain.User, name, code, parent string) (domain.Classification, error) {
	if actor.Role != domain.RoleArchivist && !actor.IsAdministrator() {
		return domain.Classification{}, domain.ErrForbidden
	}
	n := domain.Classification{ID: s.ids.New(), Name: name, Code: code, ParentID: parent, Enabled: true, CreatedAt: s.clock.Now(), UpdatedAt: s.clock.Now()}
	return n, s.repo.CreateClassification(ctx, n)
}
func (s *ClassificationService) Move(ctx context.Context, actor domain.User, id, parent string) error {
	if actor.Role != domain.RoleArchivist && !actor.IsAdministrator() {
		return domain.ErrForbidden
	}
	items, err := s.repo.ListClassifications(ctx)
	if err != nil {
		return err
	}
	nodes := map[string]domain.Classification{}
	for _, n := range items {
		nodes[n.ID] = n
	}
	if err := domain.ValidateClassificationMove(id, parent, nodes); err != nil {
		return err
	}
	node := nodes[id]
	node.ParentID = parent
	node.UpdatedAt = s.clock.Now()
	if err := s.repo.UpdateClassification(ctx, node); err != nil {
		return err
	}
	return s.repo.MoveClassification(ctx, id, parent)
}
