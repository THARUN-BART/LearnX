package execution

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrUnsupportedLanguage = errors.New("unsupported programming language")
	ErrEmptySourceCode     = errors.New("source code cannot be empty")
	ErrSourceCodeTooLarge  = errors.New("source code exceeds maximum permitted size of 64KB")
)

const MaxSourceCodeBytes = 64 * 1024 // 64KB

type LanguageConfig struct {
	ID             int    `json:"id"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	DefaultTimeout int    `json:"default_timeout_seconds"`
	DefaultMemory  int    `json:"default_memory_kb"`
}

type LanguageRegistry struct {
	mu        sync.RWMutex
	languages map[string]LanguageConfig
	aliases   map[string]string
}

var defaultRegistry *LanguageRegistry
var once sync.Once

func GetRegistry() *LanguageRegistry {
	once.Do(func() {
		defaultRegistry = &LanguageRegistry{
			languages: make(map[string]LanguageConfig),
			aliases:   make(map[string]string),
		}

		// Register standard Judge0 languages
		defaultRegistry.Register(LanguageConfig{
			ID:             50,
			Slug:           "c",
			Name:           "C",
			Version:        "GCC 9.2.0",
			DefaultTimeout: 5,
			DefaultMemory:  128000,
		})
		defaultRegistry.Register(LanguageConfig{
			ID:             54,
			Slug:           "cpp",
			Name:           "C++",
			Version:        "GCC 9.2.0",
			DefaultTimeout: 5,
			DefaultMemory:  128000,
		}, "c++", "cplusplus")
		defaultRegistry.Register(LanguageConfig{
			ID:             62,
			Slug:           "java",
			Name:           "Java",
			Version:        "OpenJDK 13.0.1",
			DefaultTimeout: 10,
			DefaultMemory:  256000,
		})
		defaultRegistry.Register(LanguageConfig{
			ID:             71,
			Slug:           "python",
			Name:           "Python",
			Version:        "3.8.1",
			DefaultTimeout: 8,
			DefaultMemory:  128000,
		}, "py", "python3")
		defaultRegistry.Register(LanguageConfig{
			ID:             60,
			Slug:           "go",
			Name:           "Go",
			Version:        "1.13.5",
			DefaultTimeout: 5,
			DefaultMemory:  128000,
		}, "golang")
	})
	return defaultRegistry
}

func (r *LanguageRegistry) Register(cfg LanguageConfig, aliases ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	canonical := strings.ToLower(strings.TrimSpace(cfg.Slug))
	r.languages[canonical] = cfg
	for _, alias := range aliases {
		r.aliases[strings.ToLower(strings.TrimSpace(alias))] = canonical
	}
}

func (r *LanguageRegistry) Resolve(lang string) (LanguageConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cleaned := strings.ToLower(strings.TrimSpace(lang))
	if target, ok := r.aliases[cleaned]; ok {
		cleaned = target
	}
	if cfg, ok := r.languages[cleaned]; ok {
		return cfg, nil
	}
	return LanguageConfig{}, ErrUnsupportedLanguage
}

func (r *LanguageRegistry) SupportedLanguages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]string, 0, len(r.languages))
	for k := range r.languages {
		list = append(list, k)
	}
	return list
}
