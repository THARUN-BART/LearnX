package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/profile"
)

func TestProfile_GetAndUpdate(t *testing.T) {
	repo := profile.NewMemoryRepository()
	svc := profile.NewService(repo)

	repo.Save(&profile.Profile{
		ID:        "student-id-1",
		Email:     "student@learnx.com",
		FullName:  "Original Name",
		College:   "Old College",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	prof, err := svc.GetProfile(context.Background(), "student-id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prof.FullName != "Original Name" {
		t.Errorf("expected Original Name, got %s", prof.FullName)
	}

	updated, err := svc.UpdateProfile(context.Background(), "student-id-1", &profile.UpdateProfileRequest{
		FullName:          "Updated Name",
		College:           "Stanford University",
		Course:            "M.Tech",
		Department:        "AI",
		PreferredLanguage: "Python",
		LearningGoal:      "DSA",
		DSALevel:          "Advanced",
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.FullName != "Updated Name" || updated.College != "Stanford University" {
		t.Errorf("unexpected updated profile: %+v", updated)
	}

	_, err = svc.GetProfile(context.Background(), "nonexistent-id")
	if !errors.Is(err, profile.ErrProfileNotFound) {
		t.Fatalf("expected ErrProfileNotFound, got %v", err)
	}
}
