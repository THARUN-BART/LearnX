package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SavePendingSignup(ctx context.Context, ps *PendingSignup) error
	GetPendingSignup(ctx context.Context, email string) (*PendingSignup, error)
	DeletePendingSignup(ctx context.Context, email string) error

	SavePasswordReset(ctx context.Context, pr *PasswordReset) error
	GetPasswordResetByTokenHash(ctx context.Context, tokenHash string) (*PasswordReset, error)
	MarkPasswordResetUsed(ctx context.Context, tokenHash string) error

	UserExists(ctx context.Context, email string) (bool, error)
	CreateProfile(ctx context.Context, p *profile.Profile) error
	GetProfileByEmail(ctx context.Context, email string) (*profile.Profile, error)
	GetProfileByID(ctx context.Context, id string) (*profile.Profile, error)
}

// PostgresRepository implements Repository using pgxpool.Pool.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const profileSelectColumns = `
	id, full_name, email, college, course, department, current_year,
	preferred_language, learning_goal, dsa_level, avatar_url,
	onboarding_completed, created_at, updated_at
`

func (r *PostgresRepository) SavePendingSignup(ctx context.Context, ps *PendingSignup) error {
	query := `
		INSERT INTO pending_signups (email, full_name, password_encrypted, password_hash, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			password_encrypted = EXCLUDED.password_encrypted,
			password_hash = EXCLUDED.password_hash,
			created_at = EXCLUDED.created_at,
			expires_at = EXCLUDED.expires_at;
	`
	_, err := r.pool.Exec(ctx, query,
		strings.ToLower(ps.Email),
		ps.FullName,
		ps.PasswordEncrypted,
		ps.PasswordHash,
		ps.CreatedAt,
		ps.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save pending signup: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetPendingSignup(ctx context.Context, email string) (*PendingSignup, error) {
	query := `
		SELECT id, email, full_name, password_encrypted, password_hash, created_at, expires_at
		FROM pending_signups
		WHERE email = $1 AND expires_at > NOW();
	`
	var ps PendingSignup
	err := r.pool.QueryRow(ctx, query, strings.ToLower(email)).Scan(
		&ps.ID,
		&ps.Email,
		&ps.FullName,
		&ps.PasswordEncrypted,
		&ps.PasswordHash,
		&ps.CreatedAt,
		&ps.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPendingSignupNotFound
		}
		return nil, fmt.Errorf("failed to get pending signup: %w", err)
	}
	return &ps, nil
}

func (r *PostgresRepository) DeletePendingSignup(ctx context.Context, email string) error {
	query := `DELETE FROM pending_signups WHERE email = $1;`
	_, err := r.pool.Exec(ctx, query, strings.ToLower(email))
	if err != nil {
		return fmt.Errorf("failed to delete pending signup: %w", err)
	}
	return nil
}

func (r *PostgresRepository) SavePasswordReset(ctx context.Context, pr *PasswordReset) error {
	query := `
		INSERT INTO password_resets (token_hash, user_id, email, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := r.pool.Exec(ctx, query,
		pr.TokenHash,
		pr.UserID,
		strings.ToLower(pr.Email),
		pr.CreatedAt,
		pr.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save password reset: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetPasswordResetByTokenHash(ctx context.Context, tokenHash string) (*PasswordReset, error) {
	query := `
		SELECT id, token_hash, user_id, email, created_at, expires_at, used_at
		FROM password_resets
		WHERE token_hash = $1;
	`
	var pr PasswordReset
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&pr.ID,
		&pr.TokenHash,
		&pr.UserID,
		&pr.Email,
		&pr.CreatedAt,
		&pr.ExpiresAt,
		&pr.UsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidResetToken
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}
	return &pr, nil
}

func (r *PostgresRepository) MarkPasswordResetUsed(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE password_resets
		SET used_at = NOW()
		WHERE token_hash = $1 AND used_at IS NULL;
	`
	cmdTag, err := r.pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to mark reset token used: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrInvalidResetToken
	}
	return nil
}

func (r *PostgresRepository) UserExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM profiles WHERE email = $1);`
	var exists bool
	err := r.pool.QueryRow(ctx, query, strings.ToLower(email)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) CreateProfile(ctx context.Context, p *profile.Profile) error {
	query := `
		INSERT INTO profiles (
			id, full_name, email, college, course, department, current_year,
			preferred_language, learning_goal, dsa_level, avatar_url, onboarding_completed, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			updated_at = NOW();
	`
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		p.ID,
		p.FullName,
		strings.ToLower(p.Email),
		p.College,
		p.Course,
		p.Department,
		p.CurrentYear,
		p.PreferredLanguage,
		p.LearningGoal,
		p.DSALevel,
		p.AvatarURL,
		p.OnboardingCompleted,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create student profile: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetProfileByEmail(ctx context.Context, email string) (*profile.Profile, error) {
	query := fmt.Sprintf("SELECT %s FROM profiles WHERE email = $1;", profileSelectColumns)
	var p profile.Profile
	err := r.pool.QueryRow(ctx, query, strings.ToLower(email)).Scan(
		&p.ID,
		&p.FullName,
		&p.Email,
		&p.College,
		&p.Course,
		&p.Department,
		&p.CurrentYear,
		&p.PreferredLanguage,
		&p.LearningGoal,
		&p.DSALevel,
		&p.AvatarURL,
		&p.OnboardingCompleted,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get profile by email: %w", err)
	}
	return &p, nil
}

func (r *PostgresRepository) GetProfileByID(ctx context.Context, id string) (*profile.Profile, error) {
	query := fmt.Sprintf("SELECT %s FROM profiles WHERE id = $1;", profileSelectColumns)
	var p profile.Profile
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.FullName,
		&p.Email,
		&p.College,
		&p.Course,
		&p.Department,
		&p.CurrentYear,
		&p.PreferredLanguage,
		&p.LearningGoal,
		&p.DSALevel,
		&p.AvatarURL,
		&p.OnboardingCompleted,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get profile by id: %w", err)
	}
	return &p, nil
}

// MemoryRepository provides a thread-safe in-memory implementation of Repository for testing.
type MemoryRepository struct {
	mu             sync.RWMutex
	pendingSignups map[string]*PendingSignup
	passwordResets map[string]*PasswordReset
	profiles       map[string]*profile.Profile
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		pendingSignups: make(map[string]*PendingSignup),
		passwordResets: make(map[string]*PasswordReset),
		profiles:       make(map[string]*profile.Profile),
	}
}

func (m *MemoryRepository) SavePendingSignup(ctx context.Context, ps *PendingSignup) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pendingSignups[strings.ToLower(ps.Email)] = ps
	return nil
}

func (m *MemoryRepository) GetPendingSignup(ctx context.Context, email string) (*PendingSignup, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ps, ok := m.pendingSignups[strings.ToLower(email)]
	if !ok || time.Now().After(ps.ExpiresAt) {
		return nil, ErrPendingSignupNotFound
	}
	return ps, nil
}

func (m *MemoryRepository) DeletePendingSignup(ctx context.Context, email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.pendingSignups, strings.ToLower(email))
	return nil
}

func (m *MemoryRepository) SavePasswordReset(ctx context.Context, pr *PasswordReset) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.passwordResets[pr.TokenHash] = pr
	return nil
}

func (m *MemoryRepository) GetPasswordResetByTokenHash(ctx context.Context, tokenHash string) (*PasswordReset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pr, ok := m.passwordResets[tokenHash]
	if !ok {
		return nil, ErrInvalidResetToken
	}
	return pr, nil
}

func (m *MemoryRepository) MarkPasswordResetUsed(ctx context.Context, tokenHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	pr, ok := m.passwordResets[tokenHash]
	if !ok || pr.UsedAt != nil {
		return ErrInvalidResetToken
	}
	now := time.Now()
	pr.UsedAt = &now
	return nil
}

func (m *MemoryRepository) UserExists(ctx context.Context, email string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.profiles {
		if strings.EqualFold(p.Email, email) {
			return true, nil
		}
	}
	return false, nil
}

func (m *MemoryRepository) CreateProfile(ctx context.Context, p *profile.Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()
	m.profiles[p.ID] = p
	return nil
}

func (m *MemoryRepository) GetProfileByEmail(ctx context.Context, email string) (*profile.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.profiles {
		if strings.EqualFold(p.Email, email) {
			return p, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MemoryRepository) GetProfileByID(ctx context.Context, id string) (*profile.Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.profiles[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return p, nil
}
