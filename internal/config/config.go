package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Provider        string            `json:"provider"`
	APIKeys         map[string]string `json:"api_keys"`
	BaseURL         string            `json:"base_url"`
	DefaultModel    string            `json:"default_model"`
	AvailableModels []string          `json:"available_models"`
	Theme           ThemeConfig       `json:"theme"`
	Context         ContextConfig     `json:"context"`
	ClearScreen     bool              `json:"clear_screen"`
	Position        string            `json:"position"`
}

type ThemeConfig struct {
	AccentColor string `json:"accent_color"`
	BorderStyle string `json:"border_style"`
}

type ContextConfig struct {
	IncludeCWD          bool `json:"include_cwd"`
	IncludeShellHistory bool `json:"include_shell_history"`
}

func DefaultConfig() Config {
	return Config{
		Provider: "openai",
		APIKeys:  map[string]string{"openai": ""},
		BaseURL:  "",
		DefaultModel:    "gpt-4o-mini",
		AvailableModels: []string{"gpt-4o-mini", "gpt-4o", "gpt-4-turbo"},
		Theme: ThemeConfig{
			AccentColor: "#FF6F61",
			BorderStyle: "rounded",
		},
		Context: ContextConfig{
			IncludeCWD:          true,
			IncludeShellHistory: false,
		},
		ClearScreen: false,
		Position:    "bottom",
	}
}

func Load() (Config, error) {
	cfg := DefaultConfig()

	path, err := ConfigFilePath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	// Check env var override
	if envKey := os.Getenv("OPENAI_API_KEY"); envKey != "" {
		if cfg.APIKeys == nil {
			cfg.APIKeys = make(map[string]string)
		}
		cfg.APIKeys["openai"] = envKey
	}

	return cfg, nil
}

func (c Config) Save() error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func (c Config) APIKey() string {
	if key, ok := c.APIKeys[c.Provider]; ok {
		return key
	}
	return ""
}
