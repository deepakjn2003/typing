package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	AppDir     = ".typing"
	DBFile     = "typing.db"
	ConfigFile = "config.json"

	TextSizeSmall  = "small"
	TextSizeMedium = "medium"
	TextSizeLarge  = "large"

	PracticeModeTimed   = "timed"
	PracticeModePassage = "passage"
)

type Config struct {
	GeminiAPIKey       string `json:"gemini_api_key,omitempty"`
	GeminiModel        string `json:"gemini_model,omitempty"`
	DBPath             string `json:"db_path,omitempty"`
	MinWords           int    `json:"min_words,omitempty"`
	MaxWords           int    `json:"max_words,omitempty"`
	TextSize           string `json:"text_size,omitempty"`
	VisibleLines       int    `json:"visible_lines,omitempty"`
	PracticeMode       string `json:"practice_mode,omitempty"`
	MaxPracticeSeconds int    `json:"max_practice_seconds,omitempty"`
}

func defaultConfig() *Config {
	return &Config{
		GeminiModel:        "gemini-2.5-flash",
		MinWords:           150,
		MaxWords:           250,
		TextSize:           TextSizeMedium,
		VisibleLines:       3,
		PracticeMode:       PracticeModePassage,
		MaxPracticeSeconds: 60,
	}
}

func Load() (*Config, error) {
	cfg := defaultConfig()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	appDir := filepath.Join(homeDir, AppDir)
	cfg.DBPath = filepath.Join(appDir, DBFile)

	configPath := filepath.Join(appDir, ConfigFile)
	if data, err := os.ReadFile(configPath); err == nil {
		var fileCfg Config
		if err := json.Unmarshal(data, &fileCfg); err == nil {
			mergeConfig(cfg, &fileCfg)
		}
	}

	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		cfg.GeminiAPIKey = key
	}
	if model := os.Getenv("TYPING_GEMINI_MODEL"); model != "" {
		cfg.GeminiModel = model
	}
	if dbPath := os.Getenv("TYPING_DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	cfg.Validate()
	return cfg, nil
}

func mergeConfig(dst, src *Config) {
	if src.GeminiAPIKey != "" {
		dst.GeminiAPIKey = src.GeminiAPIKey
	}
	if src.GeminiModel != "" {
		dst.GeminiModel = src.GeminiModel
	}
	if src.DBPath != "" {
		dst.DBPath = src.DBPath
	}
	if src.MinWords > 0 {
		dst.MinWords = src.MinWords
	}
	if src.MaxWords > 0 {
		dst.MaxWords = src.MaxWords
	}
	if src.TextSize != "" {
		dst.TextSize = src.TextSize
	}
	if src.VisibleLines > 0 {
		dst.VisibleLines = src.VisibleLines
	}
	if src.PracticeMode != "" {
		dst.PracticeMode = src.PracticeMode
	}
	if src.MaxPracticeSeconds > 0 {
		dst.MaxPracticeSeconds = src.MaxPracticeSeconds
	}
}

func (c *Config) Validate() {
	if c.TextSize != TextSizeSmall && c.TextSize != TextSizeMedium && c.TextSize != TextSizeLarge {
		c.TextSize = TextSizeMedium
	}
	if c.VisibleLines < 2 {
		c.VisibleLines = 2
	}
	if c.VisibleLines > 5 {
		c.VisibleLines = 5
	}
	if c.PracticeMode != PracticeModeTimed && c.PracticeMode != PracticeModePassage {
		c.PracticeMode = PracticeModePassage
	}
	if c.MaxWords < 50 {
		c.MaxWords = 50
	}
	if c.MaxWords > 500 {
		c.MaxWords = 500
	}
	if c.MinWords < 10 {
		c.MinWords = 10
	}
	if c.MinWords > c.MaxWords {
		c.MinWords = c.MaxWords / 2
		if c.MinWords < 10 {
			c.MinWords = 10
		}
	}
	if c.MaxPracticeSeconds < 15 {
		c.MaxPracticeSeconds = 15
	}
	if c.MaxPracticeSeconds > 600 {
		c.MaxPracticeSeconds = 600
	}
}

func (c *Config) Save() error {
	c.Validate()

	appDir, err := AppDataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(appDir, ConfigFile)
	return os.WriteFile(configPath, data, 0o644)
}

func (c *Config) IsTimedMode() bool {
	return c.PracticeMode == PracticeModeTimed
}

func AppDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppDir), nil
}

func (c *Config) HasGeminiKey() bool {
	return c.GeminiAPIKey != ""
}
