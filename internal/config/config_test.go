package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadToken_FromEnv(t *testing.T) {
	t.Setenv("LINEAR_API_KEY", "test-token-from-env")

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "test-token-from-env" {
		t.Errorf("got %q, want %q", token, "test-token-from-env")
	}
}

func TestLoadToken_FromFile(t *testing.T) {
	// Убираем переменную окружения, чтобы не мешала
	t.Setenv("LINEAR_API_KEY", "")

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte("token: file-token\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Подменяем home-директорию
	oldHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", oldHome)

	// Создаём нужную структуру директорий
	configDir := filepath.Join(dir, configDirName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, configFileName), []byte("token: file-token\n"), 0600); err != nil {
		t.Fatal(err)
	}

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "file-token" {
		t.Errorf("got %q, want %q", token, "file-token")
	}
}

func TestLoadToken_NotFound(t *testing.T) {
	t.Setenv("LINEAR_API_KEY", "")

	dir := t.TempDir()
	t.Setenv("HOME", dir)

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}

func TestSaveToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	if err := SaveToken("my-saved-token"); err != nil {
		t.Fatalf("SaveToken error: %v", err)
	}

	// Очищаем env, чтобы читать из файла
	t.Setenv("LINEAR_API_KEY", "")

	token, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken error: %v", err)
	}
	if token != "my-saved-token" {
		t.Errorf("got %q, want %q", token, "my-saved-token")
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abcdefghij", "abcd******"},
		{"abc", "***"},
		{"", ""},
		{"ab", "**"},
	}

	for _, tt := range tests {
		got := MaskToken(tt.input)
		if got != tt.want {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
