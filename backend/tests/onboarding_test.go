package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/profile"
)

func TestProfile_Onboarding_Success(t *testing.T) {
	repo := profile.NewMemoryRepository()
	svc := profile.NewService(repo)

	// Seed profile
	repo.Save(&profile.Profile{
		ID:                  "student-onboard-1",
		FullName:            "John Doe",
		Email:               "john@university.edu",
		OnboardingCompleted: false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})

	req := &profile.OnboardingRequest{
		College:           "Indian Institute of Technology",
		Course:            "B.Tech",
		Department:        "Computer Science and Engineering",
		CurrentYear:       "3rd Year",
		PreferredLanguage: "Go",
		LearningGoal:      "Interview Preparation",
		DSALevel:          "Intermediate",
	}

	p, err := svc.CompleteOnboarding(context.Background(), "student-onboard-1", req)
	if err != nil {
		t.Fatalf("onboarding failed: %v", err)
	}

	if !p.OnboardingCompleted {
		t.Errorf("expected OnboardingCompleted to be true")
	}
	if p.College != "Indian Institute of Technology" {
		t.Errorf("expected college IIT, got %s", p.College)
	}
	if p.Course != "B.Tech" {
		t.Errorf("expected course B.Tech, got %s", p.Course)
	}
	if p.Department != "Computer Science and Engineering" {
		t.Errorf("expected department CSE, got %s", p.Department)
	}
	if p.CurrentYear != "3rd Year" {
		t.Errorf("expected 3rd Year, got %s", p.CurrentYear)
	}
	if p.PreferredLanguage != "Go" {
		t.Errorf("expected preferred language Go, got %s", p.PreferredLanguage)
	}
	if p.LearningGoal != "Interview Preparation" {
		t.Errorf("expected learning goal Interview Preparation, got %s", p.LearningGoal)
	}
	if p.DSALevel != "Intermediate" {
		t.Errorf("expected DSA level Intermediate, got %s", p.DSALevel)
	}
}

func TestProfile_Onboarding_ValidationErrors(t *testing.T) {
	repo := profile.NewMemoryRepository()
	svc := profile.NewService(repo)

	repo.Save(&profile.Profile{
		ID:        "student-onboard-2",
		Email:     "student2@learnx.com",
		FullName:  "Jane Student",
		CreatedAt: time.Now(),
	})

	// 1. Invalid Preferred Language
	_, err := svc.CompleteOnboarding(context.Background(), "student-onboard-2", &profile.OnboardingRequest{
		College:           "MIT",
		Course:            "B.Sc",
		Department:        "CS",
		CurrentYear:       "2nd Year",
		PreferredLanguage: "Ruby", // Not in allowed list
		LearningGoal:      "DSA",
		DSALevel:          "Beginner",
	})
	if !errors.Is(err, profile.ErrInvalidOnboarding) {
		t.Fatalf("expected ErrInvalidOnboarding for invalid language, got: %v", err)
	}

	// 2. Invalid DSA Level
	_, err = svc.CompleteOnboarding(context.Background(), "student-onboard-2", &profile.OnboardingRequest{
		College:           "MIT",
		Course:            "B.Sc",
		Department:        "CS",
		CurrentYear:       "2nd Year",
		PreferredLanguage: "Python",
		LearningGoal:      "Competitive Programming",
		DSALevel:          "Grandmaster", // Not in allowed list
	})
	if !errors.Is(err, profile.ErrInvalidOnboarding) {
		t.Fatalf("expected ErrInvalidOnboarding for invalid DSA level, got: %v", err)
	}

	// 3. Missing College
	_, err = svc.CompleteOnboarding(context.Background(), "student-onboard-2", &profile.OnboardingRequest{
		College:           "",
		Course:            "B.Tech",
		Department:        "ECE",
		CurrentYear:       "1st Year",
		PreferredLanguage: "C++",
		LearningGoal:      "DSA",
		DSALevel:          "Beginner",
	})
	if !errors.Is(err, profile.ErrInvalidOnboarding) {
		t.Fatalf("expected ErrInvalidOnboarding for missing college, got: %v", err)
	}
}
