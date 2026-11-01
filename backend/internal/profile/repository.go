package profile

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Profile, error)
	Update(ctx context.Context, id string, req *UpdateProfileRequest) (*Profile, error)
	CompleteOnboarding(ctx context.Context, id string, req *OnboardingRequest) (*Profile, error)
}

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

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Profile, error) {
	query := fmt.Sprintf("SELECT %s FROM profiles WHERE id = $1;", profileSelectColumns)
	var p Profile
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
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to query profile: %w", err)
	}
	return &p, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, req *UpdateProfileRequest) (*Profile, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.College != "" {
		existing.College = req.College
	}
	if req.Course != "" {
		existing.Course = req.Course
	}
	if req.Department != "" {
		existing.Department = req.Department
	}
	if req.CurrentYear != "" {
		existing.CurrentYear = req.CurrentYear
	}
	if req.PreferredLanguage != "" {
		existing.PreferredLanguage = req.PreferredLanguage
	}
	if req.LearningGoal != "" {
		existing.LearningGoal = req.LearningGoal
	}
	if req.DSALevel != "" {
		existing.DSALevel = req.DSALevel
	}
	if req.AvatarURL != "" {
		existing.AvatarURL = req.AvatarURL
	}
	existing.UpdatedAt = time.Now()

	query := fmt.Sprintf(`
		UPDATE profiles
		SET full_name = $1, college = $2, course = $3, department = $4,
		    current_year = $5, preferred_language = $6, learning_goal = $7,
		    dsa_level = $8, avatar_url = $9, updated_at = $10
		WHERE id = $11
		RETURNING %s;
	`, profileSelectColumns)

	var updated Profile
	err = r.pool.QueryRow(ctx, query,
		existing.FullName,
		existing.College,
		existing.Course,
		existing.Department,
		existing.CurrentYear,
		existing.PreferredLanguage,
		existing.LearningGoal,
		existing.DSALevel,
		existing.AvatarURL,
		existing.UpdatedAt,
		id,
	).Scan(
		&updated.ID,
		&updated.FullName,
		&updated.Email,
		&updated.College,
		&updated.Course,
		&updated.Department,
		&updated.CurrentYear,
		&updated.PreferredLanguage,
		&updated.LearningGoal,
		&updated.DSALevel,
		&updated.AvatarURL,
		&updated.OnboardingCompleted,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return &updated, nil
}

func (r *PostgresRepository) CompleteOnboarding(ctx context.Context, id string, req *OnboardingRequest) (*Profile, error) {
	query := fmt.Sprintf(`
		UPDATE profiles
		SET college = $1, course = $2, department = $3, current_year = $4,
		    preferred_language = $5, learning_goal = $6, dsa_level = $7,
		    onboarding_completed = TRUE, updated_at = NOW()
		WHERE id = $8
		RETURNING %s;
	`, profileSelectColumns)

	var p Profile
	err := r.pool.QueryRow(ctx, query,
		req.College,
		req.Course,
		req.Department,
		req.CurrentYear,
		req.PreferredLanguage,
		req.LearningGoal,
		req.DSALevel,
		id,
	).Scan(
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
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to complete onboarding in database: %w", err)
	}
	return &p, nil
}

type MemoryRepository struct {
	mu       sync.RWMutex
	profiles map[string]*Profile
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		profiles: make(map[string]*Profile),
	}
}

func (m *MemoryRepository) Save(p *Profile) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.profiles[p.ID] = p
}

func (m *MemoryRepository) GetByID(ctx context.Context, id string) (*Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.profiles[id]
	if !ok {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

func (m *MemoryRepository) Update(ctx context.Context, id string, req *UpdateProfileRequest) (*Profile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.profiles[id]
	if !ok {
		return nil, ErrProfileNotFound
	}

	if req.FullName != "" {
		p.FullName = req.FullName
	}
	if req.College != "" {
		p.College = req.College
	}
	if req.Course != "" {
		p.Course = req.Course
	}
	if req.Department != "" {
		p.Department = req.Department
	}
	if req.CurrentYear != "" {
		p.CurrentYear = req.CurrentYear
	}
	if req.PreferredLanguage != "" {
		p.PreferredLanguage = req.PreferredLanguage
	}
	if req.LearningGoal != "" {
		p.LearningGoal = req.LearningGoal
	}
	if req.DSALevel != "" {
		p.DSALevel = req.DSALevel
	}
	if req.AvatarURL != "" {
		p.AvatarURL = req.AvatarURL
	}
	p.UpdatedAt = time.Now()

	return p, nil
}

func (m *MemoryRepository) CompleteOnboarding(ctx context.Context, id string, req *OnboardingRequest) (*Profile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.profiles[id]
	if !ok {
		return nil, ErrProfileNotFound
	}

	p.College = req.College
	p.Course = req.Course
	p.Department = req.Department
	p.CurrentYear = req.CurrentYear
	p.PreferredLanguage = req.PreferredLanguage
	p.LearningGoal = req.LearningGoal
	p.DSALevel = req.DSALevel
	p.OnboardingCompleted = true
	p.UpdatedAt = time.Now()

	return p, nil
}
