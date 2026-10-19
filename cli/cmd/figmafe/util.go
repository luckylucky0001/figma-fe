package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func writeJSON(path string, v any) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	b, _ := json.MarshalIndent(v, "", "  ")
	_ = os.WriteFile(path, b, 0o644)
}
