package configfile

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestReadConfigLimitsAndWrappedErrors(t *testing.T) {
	t.Chdir(t.TempDir())
	var value struct{ GridWidth int }
	if err := Read("missing.yml", &value); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing error lost its cause: %v", err)
	}
	if err := os.WriteFile("large.yml", make([]byte, (1<<20)+1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Read("large.yml", &value); err == nil || !strings.Contains(err.Error(), "1 MiB") {
		t.Fatalf("expected size limit, got %v", err)
	}
}
