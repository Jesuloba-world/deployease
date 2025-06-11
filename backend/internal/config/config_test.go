package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// set required env variables
	os.Setenv("DEPLOYEASE_JWT_SECRET", "test-jwt-secret")
	defer os.Unsetenv("DEPLOYEASE_JWT_SECRET")

	// test default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load configuration: %v", err)
	}

	// check default values
	if cfg.Server.Port != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected server host to be 0.0.0.0, got %s", cfg.Server.Host)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected database host to be localhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != "5432" {
		t.Errorf("expected database port to be 5432, got %s", cfg.Database.Port)
	}
}

func TestLoadWithEnvironmentalVariables(t *testing.T) {
	// set env variables
	os.Setenv("DEPLOYEASE_SERVER_PORT", "9090")
	os.Setenv("DEPLOYEASE_DATABASE_HOST", "testdb")
	os.Setenv("DEPLOYEASE_JWT_SECRET", "test-secret")

	// cleanup after test
	defer func() {
		os.Unsetenv("DEPLOYEASE_SERVER_PORT")
		os.Unsetenv("DEPLOYEASE_DATABASE_HOST")
		os.Unsetenv("DEPLOYEASE_JWT_SECRET")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load configuration: %v", err)
	}

	// check if env variables are loaded
	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "testdb" {
		t.Errorf("expected database host to be testdb, got %s", cfg.Database.Host)
	}

	if cfg.JWT.Secret != "test-secret" {
		t.Errorf("expected JWT secret to be test-secret, got %s", cfg.JWT.Secret)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test validation with invalid JWT secret
	os.Setenv("DEPLOYEASE_JWT_SECRET", "your-secret-key")
	defer os.Unsetenv("DEPLOYEASE_JWT_SECRET")

	_, err := Load()
	if err == nil {
		t.Error("Expected validation error for default JWT secret, but got none")
	}
}

func TestConfigurationTypes(t *testing.T) {
	// Set a valid JWT secret for testing
	os.Setenv("DEPLOYEASE_JWT_SECRET", "test-jwt-secret")
	defer os.Unsetenv("DEPLOYEASE_JWT_SECRET")

	// Test that duration fields are properly parsed
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Check that timeouts are properly parsed as durations
	expectedReadTimeout := 15 * time.Second
	if cfg.Server.ReadTimeout != expectedReadTimeout {
		t.Errorf("Expected read timeout to be %v, got %v", expectedReadTimeout, cfg.Server.ReadTimeout)
	}

	expectedWriteTimeout := 15 * time.Second
	if cfg.Server.WriteTimeout != expectedWriteTimeout {
		t.Errorf("Expected write timeout to be %v, got %v", expectedWriteTimeout, cfg.Server.WriteTimeout)
	}

	expectedIdleTimeout := 60 * time.Second
	if cfg.Server.IdleTimeout != expectedIdleTimeout {
		t.Errorf("Expected idle timeout to be %v, got %v", expectedIdleTimeout, cfg.Server.IdleTimeout)
	}

	expectedJWTExpiration := 24 * time.Hour
	if cfg.JWT.Expiration != expectedJWTExpiration {
		t.Errorf("Expected JWT expiration to be %v, got %v", expectedJWTExpiration, cfg.JWT.Expiration)
	}
}
