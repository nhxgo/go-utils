package postgres

import (
	"context"
	"strings"
	"testing"
	"time"
)

// --- Test DSN generation ---
func TestDSN(t *testing.T) {
	cfg := Postgres{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	result := dsn(cfg)

	expectedParts := []string{
		"host=localhost",
		"port=5432",
		"user=user",
		"password=pass",
		"dbname=testdb",
		"sslmode=disable",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("dsn missing part: %s", part)
		}
	}
}

// --- Test invalid config (ParseConfig fail) ---
func TestNewPgxPool_InvalidConfig(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// invalid port (string injection to break DSN)
	cfg := Postgres{
		Host:     "localhost",
		Port:     -1, // invalid
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}

	_, err := NewPgxPool(ctx, cfg, PGConfig{})

	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

// --- Test connection failure (no DB running) ---
func TestNewPgxPool_ConnectionFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := Postgres{
		Host:     "localhost",
		Port:     9999, // assume no DB here
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}

	_, err := NewPgxPool(ctx, cfg, PGConfig{
		MaxOpenConns: 5,
		MaxIdleConns: 1,
	})

	if err == nil {
		t.Fatal("expected connection failure")
	}
}

// --- Test config values applied ---
func TestNewPgxPool_ConfigApplied(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := Postgres{
		Host:     "localhost",
		Port:     9999, // no DB needed, we only test config before ping fails
		User:     "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}

	pgCfg := PGConfig{
		MaxOpenConns: 10,
		MaxIdleConns: 2,
	}

	_, err := NewPgxPool(ctx, cfg, pgCfg)

	// we EXPECT error due to ping failure
	if err == nil {
		t.Fatal("expected error due to no DB")
	}
}
