package application

import (
	"context"
	"strings"
	"time"

	"github.com/wyw14/cry002/internal/domain"
)

type AuthService struct {
	repo       Repository
	clock      Clock
	ids        IDGenerator
	passwords  PasswordHasher
	tokens     TokenManager
	refreshTTL time.Duration
}

func NewAuthService(repo Repository, clock Clock, ids IDGenerator, passwords PasswordHasher, tokens TokenManager, refreshTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, clock: clock, ids: ids, passwords: passwords, tokens: tokens, refreshTTL: refreshTTL}
}

type LoginResult struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

func (s *AuthService) Register(ctx context.Context, email, password, name string, role domain.Role) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || len(password) < 8 || name == "" {
		return domain.User{}, domain.ErrValidation
	}
	if _, err := s.repo.UserByEmail(ctx, email); err == nil {
		return domain.User{}, domain.ErrConflict
	}
	hash, err := s.passwords.Hash(password)
	if err != nil {
		return domain.User{}, err
	}
	now := s.clock.Now()
	u := domain.User{ID: s.ids.New(), Email: email, PasswordHash: hash, DisplayName: name, Role: role, Status: domain.UserActive, CreatedAt: now, UpdatedAt: now}
	return u, s.repo.CreateUser(ctx, u)
}
func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	u, err := s.repo.UserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || u.Status != domain.UserActive {
		return LoginResult{}, domain.ErrUnauthorized
	}
	if err := s.passwords.Compare(password, u.PasswordHash); err != nil {
		return LoginResult{}, domain.ErrUnauthorized
	}
	access, err := s.tokens.IssueAccess(u)
	if err != nil {
		return LoginResult{}, err
	}
	plain, hash, family, err := s.tokens.NewRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}
	rt := domain.RefreshToken{Hash: hash, UserID: u.ID, FamilyID: family, ExpiresAt: s.clock.Now().Add(s.refreshTTL)}
	if err := s.repo.StoreRefreshToken(ctx, rt); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{User: u, AccessToken: access, RefreshToken: plain}, nil
}
func (s *AuthService) Refresh(ctx context.Context, plain string) (LoginResult, error) {
	hash := s.tokens.HashRefresh(plain)
	rt, err := s.repo.ConsumeRefreshToken(ctx, hash, s.clock.Now())
	if err != nil {
		return LoginResult{}, err
	}
	u, err := s.repo.UserByID(ctx, rt.UserID)
	if err != nil {
		return LoginResult{}, err
	}
	access, err := s.tokens.IssueAccess(u)
	if err != nil {
		return LoginResult{}, err
	}
	newPlain, newHash, _, err := s.tokens.NewRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}
	if err := s.repo.StoreRefreshToken(ctx, domain.RefreshToken{Hash: newHash, UserID: u.ID, FamilyID: rt.FamilyID, ExpiresAt: s.clock.Now().Add(s.refreshTTL)}); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{User: u, AccessToken: access, RefreshToken: newPlain}, nil
}
func (s *AuthService) Logout(ctx context.Context, plain string) error {
	rt, err := s.repo.RefreshToken(ctx, s.tokens.HashRefresh(plain))
	if err != nil {
		return err
	}
	return s.repo.RevokeTokenFamily(ctx, rt.FamilyID, s.clock.Now())
}
