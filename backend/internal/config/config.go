package config

import (
	"os"
	"strconv"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/joho/godotenv"
)

type Config struct {
	OTPServiceBaseURL       string
	OTPAppName              string
	OTPTimeout              time.Duration
	OTPCooldown             time.Duration
	MaxOTPAttemptsPerWindow int

	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
	SupabaseDBURL          string
	SupabaseJWTSecret      string

	ServerPort string
	ServerEnv  string
	ClientURL  string

	EncryptionKey            []byte
	PendingSignupExpiry      time.Duration
	PasswordResetTokenExpiry time.Duration

	Judge0APIURL                string
	Judge0APIKey                string
	Judge0APIHost               string
	ExecutionTimeout            time.Duration
	SubmissionWorkerConcurrency int
	SubmissionQueueSize         int
	ExecutionProvider           string

	OpenAIAPIKey string
	OpenAIBaseURL string
	OpenAIModel  string
	GeminiAPIKey string
	GeminiBaseURL string
	GeminiModel  string
	AIProvider   string
	AITimeout    time.Duration
}

// Load loads configuration from environment variables and an optional .env file.
func Load() (*Config, error) {
	// Attempt to load .env file if it exists, ignore error if missing
	_ = godotenv.Load()

	cfg := &Config{
		OTPServiceBaseURL:       getEnv("OTP_SERVICE_BASE_URL", "https://otp-latest-0k73.onrender.com"),
		OTPAppName:              getEnv("OTP_APP_NAME", "LearnX"),
		OTPTimeout:              time.Duration(getEnvAsInt("OTP_TIMEOUT_SECONDS", 15)) * time.Second,
		OTPCooldown:             time.Duration(getEnvAsInt("OTP_COOLDOWN_SECONDS", 60)) * time.Second,
		MaxOTPAttemptsPerWindow: getEnvAsInt("OTP_MAX_ATTEMPTS", 5),

		SupabaseURL:            getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		SupabaseDBURL:          getEnv("SUPABASE_DB_URL", ""),
		SupabaseJWTSecret:      getEnv("SUPABASE_JWT_SECRET", ""),

		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerEnv:  getEnv("SERVER_ENV", "development"),
		ClientURL:  getEnv("CLIENT_URL", "http://localhost:3000"),

		PendingSignupExpiry:      time.Duration(getEnvAsInt("SIGNUP_EXPIRY_MINUTES", 10)) * time.Minute,
		PasswordResetTokenExpiry: time.Duration(getEnvAsInt("RESET_TOKEN_EXPIRY_MINUTES", 15)) * time.Minute,

		Judge0APIURL:                getEnv("JUDGE0_API_URL", "http://localhost:2358"),
		Judge0APIKey:                getEnv("JUDGE0_API_KEY", ""),
		Judge0APIHost:               getEnv("JUDGE0_API_HOST", ""),
		ExecutionTimeout:            time.Duration(getEnvAsInt("EXECUTION_TIMEOUT_SECONDS", 15)) * time.Second,
		SubmissionWorkerConcurrency: getEnvAsInt("SUBMISSION_WORKER_CONCURRENCY", 4),
		SubmissionQueueSize:         getEnvAsInt("SUBMISSION_QUEUE_SIZE", 1000),
		ExecutionProvider:           getEnv("EXECUTION_PROVIDER", "judge0"),

		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
		GeminiBaseURL: getEnv("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com/v1beta"),
		GeminiModel:   getEnv("GEMINI_MODEL", "gemini-1.5-flash"),
		AIProvider:    getEnv("AI_PROVIDER", "gemini"),
		AITimeout:     time.Duration(getEnvAsInt("AI_TIMEOUT_SECONDS", 30)) * time.Second,
	}

	rawKey := getEnv("ENCRYPTION_KEY", "")
	if rawKey == "" {
		if cfg.SupabaseServiceRoleKey != "" {
			rawKey = cfg.SupabaseServiceRoleKey
		} else {
			rawKey = "learnx-default-fallback-encryption-key-for-dev"
		}
	}
	cfg.EncryptionKey = crypto.DeriveKey(rawKey)

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
