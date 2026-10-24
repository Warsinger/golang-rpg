// Package configfile reads bounded, strict scene configuration.
package configfile

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Read accepts one strict YAML document and keeps filesystem failures wrapped.
func Read(path string, value any) error {
	return ReadFS(os.DirFS("."), path, value)
}

// ReadFS reads the same strict configuration from an injected filesystem.
func ReadFS(source fs.FS, path string, value any) error {
	f, err := source.Open(path)
	if err != nil {
		return fmt.Errorf("read %s: %w (run from the repository root)", path, err)
	}
	defer f.Close()
	const maxConfigBytes = 1 << 20
	data, err := io.ReadAll(io.LimitReader(f, maxConfigBytes+1))
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) > maxConfigBytes {
		return fmt.Errorf("%s: configuration exceeds 1 MiB", path)
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("%s: expected a single YAML document", path)
		}
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
