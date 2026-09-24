package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	logsupport "github.com/we7coreteam/w7-rangine-go/v2/pkg/support/logger"
)

func TestRotateFileChannel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	factory := NewLoggerFactory()
	factory.Register(map[string]logsupport.Config{
		"file": {Driver: "file", Path: path, Level: "debug"},
	})

	channel, err := factory.Channel("file")
	if err != nil {
		t.Fatal(err)
	}
	channel.Info("before rotation")
	if err := factory.Rotate("file"); err != nil {
		t.Fatal(err)
	}
	channel.Info("after rotation")

	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(current), "after rotation") || strings.Contains(string(current), "before rotation") {
		t.Fatalf("unexpected current log contents: %q", current)
	}
	backups, err := filepath.Glob(filepath.Join(filepath.Dir(path), "app-*.log"))
	if err != nil {
		t.Fatal(err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected one backup, got %v", backups)
	}
	old, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(old), "before rotation") || strings.Contains(string(old), "after rotation") {
		t.Fatalf("unexpected backup log contents: %q", old)
	}
}

func TestRotateRejectsOtherChannels(t *testing.T) {
	factory := NewLoggerFactory()
	factory.Register(map[string]logsupport.Config{
		"file":    {Driver: "file", Path: filepath.Join(t.TempDir(), "app.log"), Level: "info"},
		"console": {Driver: "console", Level: "info"},
		"stack":   {Driver: "stack", Channels: []string{"file", "console"}},
	})

	if err := factory.Rotate("missing"); err == nil {
		t.Fatal("expected an error for an unknown channel")
	}
	if err := factory.Rotate("file"); err == nil || !strings.Contains(err.Error(), "not initialized") {
		t.Fatalf("expected an uninitialized channel error, got %v", err)
	}
	if _, err := factory.Channel("console"); err != nil {
		t.Fatal(err)
	}
	if err := factory.Rotate("console"); err == nil || !strings.Contains(err.Error(), "not a file channel") {
		t.Fatalf("expected a non-file channel error, got %v", err)
	}
	if _, err := factory.Channel("stack"); err != nil {
		t.Fatal(err)
	}
	if err := factory.Rotate("stack"); err == nil || !strings.Contains(err.Error(), "not a file channel") {
		t.Fatalf("expected a stack channel error, got %v", err)
	}
}

func TestRotateReturnsFileError(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Join(dir, "parent")
	if err := os.WriteFile(parent, []byte("occupied"), 0600); err != nil {
		t.Fatal(err)
	}
	factory := NewLoggerFactory()
	factory.Register(map[string]logsupport.Config{
		"file": {Driver: "file", Path: filepath.Join(parent, "app.log"), Level: "info"},
	})
	if _, err := factory.Channel("file"); err != nil {
		t.Fatal(err)
	}
	if err := factory.Rotate("file"); err == nil {
		t.Fatal("expected the file rotation error")
	}
}
