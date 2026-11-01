package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidProfileData = errors.New("invalid profile update data")
	ErrInvalidOnboarding  = errors.New("invalid onboarding data")
)

type Service interface {
	GetProfile(ctx context.Context, userID string) (*Profile, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*Profile, error)
	CompleteOnboarding(ctx context.Context, userID string, req *OnboardingRequest) (*Profile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrProfileNotFound
	}
	return s.repo.GetByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*Profile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrProfileNotFound
	}
	if req.FullName != "" {
		req.FullName = strings.TrimSpace(req.FullName)
		if len(req.FullName) < 2 {
			return nil, errors.New("full name must be at least 2 characters")
		}
	}
	return s.repo.Update(ctx, userID, req)
}

func (s *service) CompleteOnboarding(ctx context.Context, userID string, req *OnboardingRequest) (*Profile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrProfileNotFound
	}

	req.College = strings.TrimSpace(req.College)
	if len(req.College) < 2 {
		return nil, fmt.Errorf("%w: college/university name is required", ErrInvalidOnboarding)
	}

	req.Course = strings.TrimSpace(req.Course)
	if req.Course == "" {
		return nil, fmt.Errorf("%w: course (e.g., B.Tech, B.E, B.Sc) is required", ErrInvalidOnboarding)
	}

	req.Department = strings.TrimSpace(req.Department)
	if req.Department == "" {
		return nil, fmt.Errorf("%w: department is required", ErrInvalidOnboarding)
	}

	req.CurrentYear = strings.TrimSpace(req.CurrentYear)
	if req.CurrentYear == "" {
		return nil, fmt.Errorf("%w: current academic year is required", ErrInvalidOnboarding)
	}

	// Validate and normalize Preferred Language
	canonicalLang, err := validatePreferredLanguage(req.PreferredLanguage)
	if err != nil {
		return nil, err
	}
	req.PreferredLanguage = canonicalLang

	// Validate Learning Goal
	req.LearningGoal = strings.TrimSpace(req.LearningGoal)
	if req.LearningGoal == "" {
		return nil, fmt.Errorf("%w: learning goal (e.g., DSA, Interview Preparation, Competitive Programming) is required", ErrInvalidOnboarding)
	}

	// Validate and normalize DSA Level
	canonicalLevel, err := validateDSALevel(req.DSALevel)
	if err != nil {
		return nil, err
	}
	req.DSALevel = canonicalLevel

	return s.repo.CompleteOnboarding(ctx, userID, req)
}

func validatePreferredLanguage(lang string) (string, error) {
	clean := strings.TrimSpace(strings.ToLower(lang))
	switch clean {
	case "c":
		return "C", nil
	case "c++", "cpp":
		return "C++", nil
	case "java":
		return "Java", nil
	case "python":
		return "Python", nil
	case "go", "golang":
		return "Go", nil
	default:
		return "", fmt.Errorf("%w: preferred language must be one of: C, C++, Java, Python, Go", ErrInvalidOnboarding)
	}
}

func validateDSALevel(level string) (string, error) {
	clean := strings.TrimSpace(strings.ToLower(level))
	switch clean {
	case "beginner":
		return "Beginner", nil
	case "intermediate":
		return "Intermediate", nil
	case "advanced":
		return "Advanced", nil
	default:
		return "", fmt.Errorf("%w: current DSA level must be one of: Beginner, Intermediate, Advanced", ErrInvalidOnboarding)
	}
}
