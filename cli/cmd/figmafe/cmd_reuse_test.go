package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChangedUIFiles(t *testing.T) {
	repo := t.TempDir()
	git := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", repo, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(rel, s string) {
		p := filepath.Join(repo, rel)
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-q")
	write("app/src/views/Old.vue", "<template/>")
	write("app/README.md", "x")
	git("add", ".")
	git("commit", "-qm", "base")

	write("app/src/views/Old.vue", "<template><div/></template>") // modified
	write("app/src/views/New.vue", "<template/>")                 // untracked
	write("app/res/layout/page.xml", "<LinearLayout/>")           // layout xml
	write("app/res/values/colors.xml", "<resources/>")            // not a UI file
	write("app/README.md", "y")                                   // not a UI file
	got, err := changedUIFiles(filepath.Join(repo, "app"), "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if l := addedLines(filepath.Join(repo, "app"), "HEAD"); !l["src/views/Old.vue"][1] || l["src/views/New.vue"] != nil {
		t.Errorf("addedLines=%v", l)
	}
	sort.Strings(got)
	want := []string{"res/layout/page.xml", "res/values/colors.xml", "src/views/New.vue", "src/views/Old.vue"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

// review: Chinese / spaced file names, renames, mode-only changes and a user's diff prefix settings
// must neither crash the gate nor silently skip files
func TestChangedFilesAndAddedLines(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	put := func(rel, text string) {
		p := filepath.Join(root, rel)
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("init", "-q")
	run("config", "diff.mnemonicPrefix", "true")
	put("src/views/订单.vue", "<template>\n<div/>\n</template>\n")
	put("src/views/My Page.vue", "<template>\n<div/>\n</template>\n")
	put("src/views/Cart.vue", "<template>\n<div/>\n</template>\n")
	put("src/views/Old.vue", "<template>\n<div/>\n</template>\n")
	run("add", "-A")
	run("commit", "-qm", "base")
	put("src/views/订单.vue", "<template>\n<div/>\n<p style=\"color:#ff0000\"/>\n</template>\n")
	put("src/views/My Page.vue", "<template>\n<div/>\n<p/>\n</template>\n")
	_ = os.Chmod(filepath.Join(root, "src/views/Cart.vue"), 0o755)
	run("mv", "src/views/Old.vue", "src/views/New.vue")
	files, err := changedUIFiles(root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(files, "|")
	for _, want := range []string{"src/views/订单.vue", "src/views/My Page.vue", "src/views/New.vue"} {
		if !strings.Contains(got, want) {
			t.Errorf("changed files %q miss %q", got, want)
		}
	}
	lines := addedLines(root, "HEAD")
	if !lines["src/views/订单.vue"][3] || len(lines["src/views/订单.vue"]) != 1 {
		t.Errorf("订单.vue added lines = %v", lines["src/views/订单.vue"])
	}
	if !lines["src/views/My Page.vue"][3] {
		t.Errorf("My Page.vue added lines = %v", lines["src/views/My Page.vue"])
	}
	if l, ok := lines["src/views/Cart.vue"]; !ok || len(l) != 0 {
		t.Errorf("a chmod-only change adds no lines, got %v (present %v)", l, ok)
	}
}

// a mistyped --focus path must fail instead of checking nothing and passing
func TestFocusPathMustExist(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"vue":"3"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runReuse([]string{"--focus", "src/views/Nope.vue", "-q", "-o", t.TempDir(), root})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("err = %v", err)
	}
}

func TestUIFileIncludesStyles(t *testing.T) {
	for f, want := range map[string]bool{"src/styles/page.scss": true, "src/a.css": true, "src/A.nvue": true, "app/src/main/res/layout/a.xml": true, "README.md": false, "src/api.ts": false} {
		if uiFile(f) != want {
			t.Errorf("uiFile(%q) = %v", f, !want)
		}
	}
}
