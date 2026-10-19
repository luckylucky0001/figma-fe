package reuse

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"
)

// A minimal markup tree for Vue templates, JSX and Android layout XML — just enough to find
// copy-pasted subtrees. Attributes and text are ignored; only the element structure counts.

type mnode struct {
	tag      string
	line     int
	children []*mnode
	size     int
	sig      string
}

var tagRe = regexp.MustCompile(`<(/?)([A-Za-z][\w.:-]*)([^<>]*?)(/?)>`)

// parseMarkup builds a forest from markup; unbalanced tags are tolerated.
func parseMarkup(src string) []*mnode {
	root := &mnode{tag: "#root"}
	stack := []*mnode{root}
	for _, m := range tagRe.FindAllStringSubmatchIndex(src, -1) {
		closing := src[m[2]:m[3]] == "/"
		tag := src[m[4]:m[5]]
		attrs := src[m[6]:m[7]]
		self := src[m[8]:m[9]] == "/" || strings.HasSuffix(strings.TrimSpace(attrs), "/")
		if strings.HasPrefix(tag, "!") || tag == "template" || tag == "script" || tag == "style" {
			continue
		}
		if closing {
			for i := len(stack) - 1; i > 0; i-- {
				if stack[i].tag == tag {
					stack = stack[:i]
					break
				}
			}
			continue
		}
		n := &mnode{tag: tag, line: strings.Count(src[:m[0]], "\n") + 1}
		parent := stack[len(stack)-1]
		parent.children = append(parent.children, n)
		if !self && !voidTag[strings.ToLower(tag)] {
			stack = append(stack, n)
		}
	}
	return root.children
}

var voidTag = map[string]bool{"img": true, "input": true, "br": true, "hr": true, "meta": true, "link": true}

// compute fills size and a shape signature. Runs of identical siblings collapse to one ("+"), so a
// copy-pasted block with 2 vs 3 rows of the same kind still counts as the same structure.
func (n *mnode) compute() {
	n.size = 1
	var parts []string
	for _, c := range n.children {
		c.compute()
		n.size += c.size
		if len(parts) > 0 && (parts[len(parts)-1] == c.sig || parts[len(parts)-1] == c.sig+"+") {
			parts[len(parts)-1] = c.sig + "+"
			continue
		}
		parts = append(parts, c.sig)
	}
	n.sig = n.tag + "(" + strings.Join(parts, ",") + ")"
}

func (n *mnode) pretty(depth, max int) string {
	if depth > max {
		return ""
	}
	s := strings.Repeat("  ", depth) + "<" + n.tag + ">\n"
	for _, c := range n.children {
		s += c.pretty(depth+1, max)
	}
	return s
}

type cloneOcc struct {
	file string
	node *mnode
}

// findClones returns subtrees with ≥minSize elements that occur ≥2 times (at least in 2 places),
// keeping only maximal clones (a clone whose parent is also a clone at every location is dropped).
func findClones(files map[string]string, parse func(string) []*mnode, minSize int) []Clone {
	occ := map[string][]cloneOcc{}
	parent := map[*mnode]*mnode{}
	var walk func(f string, n *mnode)
	walk = func(f string, n *mnode) {
		if n.size >= minSize {
			occ[n.sig] = append(occ[n.sig], cloneOcc{f, n})
		}
		for _, c := range n.children {
			parent[c] = n
			walk(f, c)
		}
	}
	names := make([]string, 0, len(files))
	for f := range files {
		names = append(names, f)
	}
	sort.Strings(names)
	for _, f := range names {
		for _, n := range parse(files[f]) {
			n.compute()
			walk(f, n)
		}
	}
	var out []Clone
	for sig, os := range occ {
		if len(os) < 2 {
			continue
		}
		maximal := false
		for _, o := range os {
			p := parent[o.node]
			if p == nil || len(occ[p.sig]) < len(os) {
				maximal = true
				break
			}
		}
		if !maximal {
			continue
		}
		c := Clone{Signature: short(sig), Elements: os[0].node.size, Sample: os[0].node.pretty(0, 3)}
		for _, o := range os {
			c.Locations = append(c.Locations, o.file+":"+itoa(o.node.line))
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		si, sj := out[i].Elements*len(out[i].Locations), out[j].Elements*len(out[j].Locations)
		if si != sj {
			return si > sj
		}
		return out[i].Signature < out[j].Signature
	})
	return out
}

func short(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:4])
}
