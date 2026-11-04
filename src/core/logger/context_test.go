package logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	logsupport "github.com/we7coreteam/w7-rangine-go/v2/pkg/support/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func logContextFixture(t *testing.T, channel string) (*zap.Logger, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "context.log")
	factory := NewLoggerFactory()
	factory.Register(map[string]logsupport.Config{
		"file":   {Driver: "file", Path: path, Level: "debug"},
		"stack":  {Driver: "stack", Channels: []string{"file"}},
		"nested": {Driver: "stack", Channels: []string{"stack"}},
	})
	// Resolve the outputs before concurrent writes; this test targets immutable
	// logger context, not factory channel initialization.
	for _, name := range []string{"file", "stack", "nested"} {
		if _, err := factory.Channel(name); err != nil {
			t.Fatal(err)
		}
	}
	log, err := factory.Channel(channel)
	if err != nil {
		t.Fatal(err)
	}
	return log, path
}

func readLogFields(t *testing.T, path string) []map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var records []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		start := strings.Index(line, "{")
		if start < 0 {
			t.Fatalf("missing structured fields: %s", line)
		}
		var fields map[string]any
		if err := json.Unmarshal([]byte(line[start:]), &fields); err != nil {
			t.Fatal(err)
		}
		records = append(records, fields)
	}
	return records
}

func TestSlogWithContextAcrossChannels(t *testing.T) {
	for _, channel := range []string{"file", "stack", "nested"} {
		t.Run(channel, func(t *testing.T) {
			output, path := logContextFixture(t, channel)
			parent := slog.New(zapslog.NewHandler(output.Core())).With("resource_type", "bot")
			child := parent.WithGroup("order").With("order_sn", "local-1")
			child.Info("payment confirmation already applied", "external_order_sn", "external-1")
			parent.Info("parent remains unchanged", "root", true)
			records := readLogFields(t, path)
			if len(records) != 2 {
				t.Fatalf("got %d records", len(records))
			}
			order, ok := records[0]["order"].(map[string]any)
			if !ok || order["order_sn"] != "local-1" || order["external_order_sn"] != "external-1" || records[0]["resource_type"] != "bot" {
				t.Fatalf("lost context: %+v", records[0])
			}
			if _, leaked := records[1]["order"]; leaked || records[1]["resource_type"] != "bot" || records[1]["root"] != true {
				t.Fatalf("child context leaked: %+v", records[1])
			}
		})
	}
}

func TestZapWithContextIsolatedDuringConcurrentWrites(t *testing.T) {
	output, path := logContextFixture(t, "nested")
	parent := output.With(zap.String("resource_type", "bot"))
	var wg sync.WaitGroup
	const count = 32
	for i := 0; i < count; i++ {
		child := parent.With(zap.String("order_sn", fmt.Sprint(i)))
		wg.Add(1)
		go func(id int, child *zap.Logger) { defer wg.Done(); child.Info("order event", zap.Int("id", id)) }(i, child)
	}
	wg.Wait()
	records := readLogFields(t, path)
	if len(records) != count {
		t.Fatalf("got %d records", len(records))
	}
	for _, record := range records {
		if record["resource_type"] != "bot" || record["order_sn"] != fmt.Sprint(int(record["id"].(float64))) {
			t.Fatalf("context leaked or missing: %+v", record)
		}
	}
}
