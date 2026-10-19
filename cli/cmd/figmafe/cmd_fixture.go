package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luckylucky0001/figma-fe/cli/internal/config"
	"github.com/luckylucky0001/figma-fe/cli/internal/fixture"
)

// parse lets flags appear before or after positional args.
func parse(fs *flag.FlagSet, args []string) []string {
	var pos []string
	for len(args) > 0 {
		_ = fs.Parse(args)
		args = fs.Args()
		if len(args) > 0 {
			pos = append(pos, args[0])
			args = args[1:]
		}
	}
	return pos
}

func load(root string) (*config.Config, error) {
	if root == "" {
		root, _ = os.Getwd()
	}
	return config.Load(root)
}

func runFixture(args []string) error {
	fs := flag.NewFlagSet("fixture", flag.ExitOnError)
	out := fs.String("o", "", "output export dir")
	root := fs.String("root", "", "repo root")
	pos := parse(fs, args)
	if len(pos) != 1 || *out == "" {
		return fmt.Errorf("usage: figmafe fixture <design.html> -o <exportDir>")
	}
	cfg, err := load(*root)
	if err != nil {
		return err
	}
	srcLib := filepath.Join(filepath.Dir(pos[0]), "library.json")
	if err := fixture.Build(pos[0], srcLib, *out); err != nil {
		return err
	}
	// Publish the library listing next to the page exports, as the plugin does for a library file.
	b, err := os.ReadFile(srcLib)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfg.P(cfg.Figma.Library), b, 0o644); err != nil {
		return err
	}
	fmt.Printf("✓ wrote %s/export.json, frame@2x.png and %s\n", *out, cfg.Figma.Library)
	return nil
}
