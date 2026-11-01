package profile

import "time"

// Profile represents the full student profile stored in Supabase PostgreSQL.
type Profile struct {
	ID                  string    `json:"id"`
	FullName            string    `json:"full_name"`
	Email               string    `json:"email"`
	College             string    `json:"college"`
	Course              string    `json:"course"`
	Department          string    `json:"department"`
	CurrentYear         string    `json:"current_year"`
	PreferredLanguage   string    `json:"preferred_language"`
	LearningGoal        string    `json:"learning_goal"`
	DSALevel            string    `json:"dsa_level"`
	AvatarURL           string    `json:"avatar_url"`
	OnboardingCompleted bool      `json:"onboarding_completed"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// OnboardingRequest contains the required student onboarding fields.
type OnboardingRequest struct {
	College           string `json:"college"`
	Course            string `json:"course"`
	Department        string `json:"department"`
	CurrentYear       string `json:"current_year"`
	PreferredLanguage string `json:"preferred_language"`
	LearningGoal      string `json:"learning_goal"`
	DSALevel          string `json:"dsa_level"`
}

// UpdateProfileRequest allows updating profile details.
type UpdateProfileRequest struct {
	FullName          string `json:"full_name,omitempty"`
	College           string `json:"college,omitempty"`
	Course            string `json:"course,omitempty"`
	Department        string `json:"department,omitempty"`
	CurrentYear       string `json:"current_year,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	LearningGoal      string `json:"learning_goal,omitempty"`
	DSALevel          string `json:"dsa_level,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
}
