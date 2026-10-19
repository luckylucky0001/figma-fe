// Package config loads figmafe.yaml from the repo root.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Root  string `yaml:"-"`
	Figma struct {
		Library   string            `yaml:"library"`
		Variables string            `yaml:"variables"`
		Pages     map[string]string `yaml:"pages"`
	} `yaml:"figma"`
	Map string `yaml:"map"`
	Out string `yaml:"out"`
	H5  struct {
		Root  string `yaml:"root"`
		DS    string `yaml:"ds"`
		Pages string `yaml:"pages"`
		Build string `yaml:"build"`
		Dist  string `yaml:"dist"`
	} `yaml:"h5"`
	Android struct {
		Root    string `yaml:"root"`
		DS      string `yaml:"ds"`
		DSView  string `yaml:"dsView"`
		App     string `yaml:"app"`
		Screens string `yaml:"screens"`
		Build   string `yaml:"build"`
		Run     string `yaml:"run"`
	} `yaml:"android"`
	Thresholds Thresholds `yaml:"thresholds"`
}

type Thresholds struct {
	Reuse         float64            `yaml:"reuse"`
	Violations    int                `yaml:"violations"`
	PixelMismatch map[string]float64 `yaml:"pixelMismatch"`
	NodeErrors    int                `yaml:"nodeErrors"`
	PositionTol   float64            `yaml:"positionTolerance"`
}

// Find walks up from dir until it finds figmafe.yaml.
func Find(dir string) (string, error) {
	d, _ := filepath.Abs(dir)
	for {
		if _, err := os.Stat(filepath.Join(d, "figmafe.yaml")); err == nil {
			return d, nil
		}
		p := filepath.Dir(d)
		if p == d {
			return "", fmt.Errorf("figmafe.yaml not found from %s upward", dir)
		}
		d = p
	}
}

func Load(root string) (*Config, error) {
	r, err := Find(root)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(filepath.Join(r, "figmafe.yaml"))
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("figmafe.yaml: %w", err)
	}
	c.Root = r
	if c.Out == "" {
		c.Out = ".figmafe"
	}
	if c.Thresholds.Reuse == 0 {
		c.Thresholds.Reuse = 0.95
	}
	if c.Thresholds.PositionTol == 0 {
		c.Thresholds.PositionTol = 2
	}
	if c.Thresholds.PixelMismatch == nil {
		c.Thresholds.PixelMismatch = map[string]float64{"h5": 0.03, "android": 0.04}
	}
	return &c, nil
}

// P joins a repo-relative path onto the root.
func (c *Config) P(parts ...string) string {
	return filepath.Join(append([]string{c.Root}, parts...)...)
}

// PageDir returns the Figma export dir of a page.
func (c *Config) PageDir(page string) (string, error) {
	d, ok := c.Figma.Pages[page]
	if !ok {
		return "", fmt.Errorf("page %q not in figmafe.yaml figma.pages", page)
	}
	return c.P(d), nil
}

// OutDir returns .figmafe/<platform>/<page>, creating it.
func (c *Config) OutDir(platform, page string) string {
	d := c.P(c.Out, platform, page)
	_ = os.MkdirAll(d, 0o755)
	return d
}
