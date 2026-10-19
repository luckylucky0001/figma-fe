package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
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
