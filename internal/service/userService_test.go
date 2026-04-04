package service_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/service"
)

func newTestUserService(userRepo *MockUserRepository, profileRepo *MockProfileRepository) service.UserService {
	return service.NewUserService(userRepo, profileRepo)
}

// ─── ChangePassword ──────────────────────────────────────────────────────────

func TestUserService_ChangePassword(t *testing.T) {
	ctx := context.Background()
	correctHash, _ := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.MinCost)

	tests := []struct {
		name        string
		userID      int
		currentPwd  string
		newPwd      string
		userRepo    *MockUserRepository
		profileRepo *MockProfileRepository
		wantErr     bool
	}{
		{
			name:       "success - password changed",
			userID:     1,
			currentPwd: "old-password",
			newPwd:     "new-password",
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, PasswordHash: string(correctHash)}, nil
				},
				UpdateFn: func(_ context.Context, _ *domain.User) error { return nil },
			},
			profileRepo: &MockProfileRepository{},
			wantErr:     false,
		},
		{
			name:       "error - wrong current password",
			userID:     1,
			currentPwd: "wrong-password",
			newPwd:     "new-password",
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, PasswordHash: string(correctHash)}, nil
				},
			},
			profileRepo: &MockProfileRepository{},
			wantErr:     true,
		},
		{
			name:       "error - user not found",
			userID:     999,
			currentPwd: "any",
			newPwd:     "new",
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return nil, domain.ErrNotFound
				},
			},
			profileRepo: &MockProfileRepository{},
			wantErr:     true,
		},
		{
			name:       "error - update failure",
			userID:     1,
			currentPwd: "old-password",
			newPwd:     "new-password",
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, PasswordHash: string(correctHash)}, nil
				},
				UpdateFn: func(_ context.Context, _ *domain.User) error {
					return errors.New("db write error")
				},
			},
			profileRepo: &MockProfileRepository{},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := newTestUserService(tc.userRepo, tc.profileRepo)
			err := svc.ChangePassword(ctx, tc.userID, tc.currentPwd, tc.newPwd)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// New password must be hashed (not stored as plain text).
func TestUserService_ChangePassword_NewPasswordIsHashed(t *testing.T) {
	ctx := context.Background()
	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old"), bcrypt.MinCost)

	var savedHash string
	svc := newTestUserService(
		&MockUserRepository{
			GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
				return &domain.User{ID: 1, PasswordHash: string(oldHash)}, nil
			},
			UpdateFn: func(_ context.Context, u *domain.User) error {
				savedHash = u.PasswordHash
				return nil
			},
		},
		&MockProfileRepository{},
	)

	if err := svc.ChangePassword(ctx, 1, "old", "new-secure-pass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedHash == "new-secure-pass" {
		t.Error("password stored as plain text - must be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte("new-secure-pass")); err != nil {
		t.Errorf("saved hash does not match new password: %v", err)
	}
}

// ─── DeleteAccount ───────────────────────────────────────────────────────────

func TestUserService_DeleteAccount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      int
		userRepo    *MockUserRepository
		wantErr     bool
		wantStatus  string
	}{
		{
			name:   "success - account soft-deleted",
			userID: 1,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, Status: "active"}, nil
				},
				UpdateFn: func(_ context.Context, _ *domain.User) error { return nil },
			},
			wantErr:    false,
			wantStatus: "deleted",
		},
		{
			name:   "error - user not found",
			userID: 999,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return nil, domain.ErrNotFound
				},
			},
			wantErr: true,
		},
		{
			name:   "error - update failure",
			userID: 1,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, Status: "active"}, nil
				},
				UpdateFn: func(_ context.Context, _ *domain.User) error {
					return errors.New("write failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var capturedUser *domain.User
			if tc.userRepo.UpdateFn != nil {
				orig := tc.userRepo.UpdateFn
				tc.userRepo.UpdateFn = func(ctx context.Context, u *domain.User) error {
					capturedUser = u
					return orig(ctx, u)
				}
			}

			svc := newTestUserService(tc.userRepo, &MockProfileRepository{})
			err := svc.DeleteAccount(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantStatus != "" && capturedUser != nil && capturedUser.Status != tc.wantStatus {
				t.Errorf("status = %q, want %q", capturedUser.Status, tc.wantStatus)
			}
		})
	}
}

// ─── GetProfile ──────────────────────────────────────────────────────────────

func TestUserService_GetProfile(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      int
		userRepo    *MockUserRepository
		profileRepo *MockProfileRepository
		wantUser    bool
		wantProfile bool
		wantErr     bool
	}{
		{
			name:   "success - user with existing profile",
			userID: 1,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1, Email: "u@example.com"}, nil
				},
			},
			profileRepo: &MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return &domain.Profile{ID: 1, UserID: 1, FirstName: "Alice"}, nil
				},
			},
			wantUser:    true,
			wantProfile: true,
		},
		{
			name:   "success - user without profile returns nil profile (not error)",
			userID: 2,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 2, Email: "new@example.com"}, nil
				},
			},
			profileRepo: &MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return nil, domain.ErrNotFound
				},
			},
			wantUser:    true,
			wantProfile: false, // nil profile is expected
		},
		{
			name:   "error - user not found",
			userID: 999,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return nil, domain.ErrNotFound
				},
			},
			profileRepo: &MockProfileRepository{},
			wantErr:     true,
		},
		{
			name:   "error - profile repo unexpected error",
			userID: 1,
			userRepo: &MockUserRepository{
				GetByIDFn: func(_ context.Context, _ int) (*domain.User, error) {
					return &domain.User{ID: 1}, nil
				},
			},
			profileRepo: &MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return nil, errors.New("db connection lost")
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := newTestUserService(tc.userRepo, tc.profileRepo)
			user, profile, err := svc.GetProfile(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantUser && user == nil {
				t.Error("expected user, got nil")
			}
			if tc.wantProfile && profile == nil {
				t.Error("expected profile, got nil")
			}
			if !tc.wantProfile && profile != nil {
				t.Errorf("expected nil profile, got %+v", profile)
			}
		})
	}
}

// ─── UpdateProfile ───────────────────────────────────────────────────────────

func TestUserService_UpdateProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("success - creates profile when not existing", func(t *testing.T) {
		var created bool
		svc := newTestUserService(
			&MockUserRepository{},
			&MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return nil, domain.ErrNotFound
				},
				CreateFn: func(_ context.Context, _ *domain.Profile) error {
					created = true
					return nil
				},
			},
		)
		_, err := svc.UpdateProfile(ctx, 1, &domain.Profile{FirstName: "Bob", LastName: "Smith"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !created {
			t.Error("expected profile to be created")
		}
	})

	t.Run("success - partial update of existing profile", func(t *testing.T) {
		existing := &domain.Profile{ID: 1, UserID: 1, FirstName: "Alice", LastName: "Smith"}
		var updated bool
		svc := newTestUserService(
			&MockUserRepository{},
			&MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return existing, nil
				},
				UpdateFn: func(_ context.Context, p *domain.Profile) error {
					updated = true
					if p.FirstName != "Bob" {
						return errors.New("FirstName should have been updated to Bob")
					}
					if p.LastName != "Smith" {
						return errors.New("LastName should remain Smith")
					}
					return nil
				},
			},
		)
		_, err := svc.UpdateProfile(ctx, 1, &domain.Profile{FirstName: "Bob"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !updated {
			t.Error("expected profile.Update to be called")
		}
	})

	t.Run("error - profile repo failure on Get", func(t *testing.T) {
		svc := newTestUserService(
			&MockUserRepository{},
			&MockProfileRepository{
				GetByUserIDFn: func(_ context.Context, _ int) (*domain.Profile, error) {
					return nil, errors.New("db error")
				},
			},
		)
		_, err := svc.UpdateProfile(ctx, 1, &domain.Profile{FirstName: "Bob"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
