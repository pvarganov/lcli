package cmd

import (
	"testing"
)

func TestRootCmdInit(t *testing.T) {
	if rootCmd.Use != "lcli" {
		t.Errorf("expected Use=lcli, got %s", rootCmd.Use)
	}
	if rootCmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}

func TestRootCmdHasTokenFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("expected --token flag to be registered")
	}
	if flag.DefValue != "" {
		t.Errorf("expected default empty token, got %s", flag.DefValue)
	}
}

func TestGetToken(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()
	if GetToken() != "test-token" {
		t.Errorf("expected test-token, got %s", GetToken())
	}
}
