package capture

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Node is one rendered element (1x units), shared by both platforms (docs/CONTRACTS.md §4).
type Node struct {
	FigmaID    *string `json:"figmaId"`
	DS         *string `json:"ds"`
	Text       *string `json:"text"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	W          float64 `json:"w"`
	H          float64 `json:"h"`
	FontSize   float64 `json:"fontSize,omitempty"`
	FontWeight float64 `json:"fontWeight,omitempty"`
	Color      string  `json:"color,omitempty"`
}

type Nodes struct {
	Platform string  `json:"platform"`
	Scale    float64 `json:"scale"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Nodes    []Node  `json:"nodes"`
}

func LoadNodes(path string) (*Nodes, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var n Nodes
	return &n, json.Unmarshal(b, &n)
}

const h5NodesJS = `(() => {
  const out = [];
  for (const el of document.querySelectorAll('body *')) {
    const r = el.getBoundingClientRect();
    if (r.width === 0 && r.height === 0) continue;
    const fid = el.getAttribute('data-figma-id');
    const ds = el.getAttribute('data-ds');
    let own = '';
    for (const n of el.childNodes) if (n.nodeType === 3) own += n.textContent;
    own = own.replace(/\s+/g, ' ').trim();
    if (!fid && !ds && !own) continue;
    const cs = getComputedStyle(el);
    out.push({ figmaId: fid, ds, text: own || null, x: r.left, y: r.top, w: r.width, h: r.height,
      fontSize: parseFloat(cs.fontSize), fontWeight: parseFloat(cs.fontWeight), color: cs.color });
  }
  return JSON.stringify(out);
})()`

// H5 builds (optional), serves dist, and captures page at w×h CSS px @2x.
func H5(h5Root, buildCmd, dist, page, outDir string, w, h int, build bool) error {
	if build && buildCmd != "" {
		cmd := exec.Command("sh", "-c", buildCmd)
		cmd.Dir = h5Root
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("h5 build failed: %v\n%s", err, tail(string(out), 40))
		}
	}
	distDir := filepath.Join(h5Root, dist)
	if _, err := os.Stat(filepath.Join(distDir, "index.html")); err != nil {
		return fmt.Errorf("%s/index.html missing (build first)", distDir)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	srv := &http.Server{Handler: http.FileServer(http.Dir(distDir))}
	go func() { _ = srv.Serve(ln) }()
	defer srv.Close()
	url := fmt.Sprintf("http://%s/index.html?page=%s", ln.Addr().String(), page)

	b, err := NewBrowser(90 * time.Second)
	if err != nil {
		return err
	}
	defer b.Close()
	if err := b.Open(url, w, h, 2); err != nil {
		return err
	}
	if err := b.Screenshot(filepath.Join(outDir, "screen.png"), w, h, 2); err != nil {
		return err
	}
	var raw string
	if err := b.Eval(h5NodesJS, &raw); err != nil {
		return err
	}
	var nodes []Node
	if err := json.Unmarshal([]byte(raw), &nodes); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("page %q rendered nothing (is src/pages/%s/index.tsx present?)", page, page)
	}
	return writeNodes(filepath.Join(outDir, "nodes.json"), &Nodes{Platform: "h5", Scale: 2, Width: float64(w), Height: float64(h), Nodes: nodes})
}

// Android builds the harness (optional) and renders the registered screen.
func Android(aRoot, buildCmd, runCmd, page, outDir string, w, h int, build bool) error {
	if build && buildCmd != "" {
		cmd := exec.Command("sh", "-c", buildCmd)
		cmd.Dir = aRoot
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("android build failed: %v\n%s", err, tail(string(out), 40))
		}
	}
	cmd := exec.Command(filepath.Join(aRoot, runCmd), page, outDir, fmt.Sprint(w), fmt.Sprint(h))
	cmd.Dir = aRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("android render failed: %v\n%s", err, tail(string(out), 40))
	}
	return nil
}

func writeNodes(path string, n *Nodes) error {
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
