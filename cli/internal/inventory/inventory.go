// Package inventory scans component assets (Phase 0) into a manifest the Agent can read
// instead of the raw source tree: components, props/variant axes, examples, usage counts, tokens.
package inventory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Prop struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Values   []string `json:"values,omitempty"` // variant values for literal unions / enums
	Optional bool     `json:"optional,omitempty"`
	Default  string   `json:"default,omitempty"`
	Doc      string   `json:"doc,omitempty"`
}

type Component struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"` // react | compose | view
	File        string   `json:"file"`
	Import      string   `json:"import"`
	Description string   `json:"description,omitempty"`
	Example     string   `json:"example,omitempty"`
	Props       []Prop   `json:"props"`
	Usage       int      `json:"usage"`
	UsageFiles  []string `json:"usageFiles,omitempty"`
}

type Token struct {
	Name   string `json:"name"` // canonical, e.g. color.text-primary
	Value  string `json:"value"`
	Source string `json:"source"`
}

// Finding is a pre-existing bad practice found in business code (adoption baseline).
type Finding struct {
	Rule string `json:"rule"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

type Inventory struct {
	Schema     string      `json:"schema"`
	Platform   string      `json:"platform"`
	Root       string      `json:"root"`
	Components []Component `json:"components"`
	Tokens     []Token     `json:"tokens"`
	Baseline   []Finding   `json:"baseline"`
}

func (inv *Inventory) Find(name string) *Component {
	for i := range inv.Components {
		if inv.Components[i].Name == name {
			return &inv.Components[i]
		}
	}
	return nil
}

func (inv *Inventory) Save(path string) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	b, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func Load(path string) (*Inventory, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inv Inventory
	return &inv, json.Unmarshal(b, &inv)
}

// ---- shared text helpers ----

func readFile(p string) string {
	b, _ := os.ReadFile(p)
	return string(b)
}

// matchClose returns the index of the bracket closing the one at open (s[open] is '(' '{' or '[').
func matchClose(s string, open int) int {
	pairs := map[byte]byte{'(': ')', '{': '}', '[': ']'}
	want := []byte{pairs[s[open]]}
	inStr := byte(0)
	for i := open + 1; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			if c == '\\' {
				i++
			} else if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			inStr = c
		case '(', '{', '[':
			want = append(want, pairs[c])
		case ')', '}', ']':
			if len(want) > 0 && c == want[len(want)-1] {
				want = want[:len(want)-1]
				if len(want) == 0 {
					return i
				}
			}
		}
	}
	return -1
}

// splitTop splits on sep at nesting depth 0.
func splitTop(s string, sep byte) []string {
	var out []string
	depth, start := 0, 0
	inStr := byte(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = c
		case '(', '{', '[', '<':
			depth++
		case ')', '}', ']', '>':
			if !(c == '>' && i > 0 && s[i-1] == '-') {
				depth--
			}
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	if strings.TrimSpace(s[start:]) != "" {
		out = append(out, s[start:])
	}
	return out
}

var docRe = regexp.MustCompile(`(?s)/\*\*(.*?)\*/\s*$`)

// docBefore returns the cleaned /** */ block that ends right before idx (ignoring annotations).
func docBefore(src string, idx int) (desc, example string) {
	head := src[:idx]
	head = strings.TrimRight(head, " \t\r\n")
	for strings.HasSuffix(head, ")") || regexp.MustCompile(`@\w+\s*$`).MatchString(head) {
		// strip a trailing annotation like @Composable or @JvmOverloads
		m := regexp.MustCompile(`@\w+(\([^)]*\))?\s*$`).FindStringIndex(head)
		if m == nil {
			break
		}
		head = strings.TrimRight(head[:m[0]], " \t\r\n")
	}
	// Only the last /** */ block directly before idx belongs to the declaration (earlier ones document props).
	if i := strings.LastIndex(head, "/**"); i >= 0 {
		head = head[i:]
	}
	m := docRe.FindStringSubmatch(head)
	if m == nil {
		return "", ""
	}
	var lines []string
	for _, l := range strings.Split(m[1], "\n") {
		l = strings.TrimSpace(l)
		l = strings.TrimPrefix(l, "*")
		if strings.HasPrefix(l, " ") {
			l = l[1:]
		}
		lines = append(lines, l)
	}
	var d, ex []string
	inEx := false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "@example") || strings.HasPrefix(t, "@sample") {
			inEx = true
			rest := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(t, "@example"), "@sample"))
			if rest != "" {
				ex = append(ex, rest)
			}
			continue
		}
		if strings.HasPrefix(t, "@") {
			inEx = false
			continue
		}
		if t == "```" || strings.HasPrefix(t, "```") {
			continue
		}
		if inEx {
			if t != "" {
				ex = append(ex, l)
			}
		} else if t != "" {
			d = append(d, t)
		}
	}
	return strings.Join(d, " "), strings.Join(ex, "\n")
}

func lineOf(src string, idx int) int { return strings.Count(src[:idx], "\n") + 1 }

func walkFiles(root string, exts []string, skip func(string) bool) []string {
	var out []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			n := info.Name()
			if n == "node_modules" || n == "build" || n == "dist" || n == ".gradle" || strings.HasPrefix(n, ".") {
				return filepath.SkipDir
			}
			if skip != nil && skip(p) {
				return filepath.SkipDir
			}
			return nil
		}
		for _, e := range exts {
			if strings.HasSuffix(p, e) {
				out = append(out, p)
				break
			}
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// kebab converts CamelCase to kebab-case: TextPrimary → text-primary.
func kebab(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var hexRe = regexp.MustCompile(`#[0-9a-fA-F]{6}\b|#[0-9a-fA-F]{3}\b`)

func rel(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return r
}
