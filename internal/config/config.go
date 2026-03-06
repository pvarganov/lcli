package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const configDirName = ".config/lcli"
const configFileName = "config.yaml"

// Config хранит настройки lcli.
type Config struct {
	Token string `yaml:"token"`
}

// LoadToken читает токен из переменной окружения LINEAR_API_KEY или из файла конфига.
// Возвращает пустую строку, если токен не найден.
func LoadToken() (string, error) {
	if token := strings.TrimSpace(os.Getenv("LINEAR_API_KEY")); token != "" {
		return token, nil
	}

	path, err := ConfigFilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parse config: %w", err)
	}

	return strings.TrimSpace(cfg.Token), nil
}

// SaveToken сохраняет токен в файл конфига.
func SaveToken(token string) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	cfg := Config{Token: token}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// ConfigFilePath возвращает абсолютный путь к файлу конфига.
func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, configDirName, configFileName), nil
}

// MaskToken возвращает замаскированную версию токена (первые 4 символа + ***).
func MaskToken(token string) string {
	if len(token) <= 4 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-4)
}
