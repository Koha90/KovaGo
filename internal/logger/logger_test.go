package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Koha90/KovaGo/internal/config"
)

func TestNewWritesLogToFile(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"kovago.log",
	)

	cfg := config.Log{
		Level:  "debug",
		Format: "text",
		File: config.LogFile{
			Enabled: true,
			Path:    path,
			Format:  "json",
		},
	}

	log, closeLogger, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	log.Info(
		"test message",
		"value",
		42,
	)

	if err := closeLogger(); err != nil {
		t.Fatalf("close logger: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	var record map[string]any

	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf(
			"unmarshal log record: %v",
			err,
		)
	}

	if got := record["msg"]; got != "test message" {
		t.Errorf(
			"msg = %v, want %q",
			got,
			"test message",
		)
	}

	if got := record["value"]; got != float64(42) {
		t.Errorf(
			"value = %v, want %v",
			got,
			42,
		)
	}
}
