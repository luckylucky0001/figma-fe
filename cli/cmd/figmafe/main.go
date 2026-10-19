// Command figmafe: Figma → H5 / Android component-reuse toolkit.
//
//	figmafe fixture   <design.html> -o <exportDir>          build a Figma-format export from an annotated design HTML
//	figmafe scan      --platform h5|android                  Phase 0: scan component assets → inventory
//	figmafe map       suggest|check                           Phase 1: Figma ↔ code mapping
//	figmafe catalog                                           generate Skill catalog / rules from inventory + map
//	figmafe plan      --page P --platform h5|android          Phase 2: reuse plan for a Figma page
//	figmafe audit     --page P --platform h5|android          static audit of the implemented page
//	figmafe capture   --page P --platform h5|android          render page → screen.png + nodes.json
//	figmafe diff      --page P --platform h5|android          pixel + node-level diff vs Figma
//	figmafe verify    --page P --platform h5|android          audit + capture + diff + report (the Skill's loop step)
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "fixture":
		err = runFixture(args)
	case "scan":
		err = runScan(args)
	case "map":
		err = runMap(args)
	case "catalog":
		err = runCatalog(args)
	case "plan":
		err = runPlan(args)
	case "audit":
		err = runAudit(args)
	case "capture":
		err = runCapture(args)
	case "diff":
		err = runDiff(args)
	case "verify":
		err = runVerify(args)
	case "reuse":
		err = runReuse(args)
	case "compare":
		err = runCompare(args)
	case "match":
		err = runMatch(args)
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", cmd)
		usage()
		os.Exit(2)
	}
	if err != nil {
		if e, ok := err.(exitError); ok {
			fmt.Fprintln(os.Stderr, e.msg)
			os.Exit(e.code)
		}
		fmt.Fprintln(os.Stderr, "figmafe:", err)
		os.Exit(1)
	}
}

type exitError struct {
	code int
	msg  string
}

func (e exitError) Error() string { return e.msg }

func usage() {
	fmt.Fprint(os.Stderr, `figmafe — Figma → H5 / Android component reuse toolkit

Phase 0  scan     --platform h5|android        scan component assets (inventory.json + catalog)
Phase 1  map      suggest | check              suggest / check Figma ↔ code mapping (gap + token drift report)
         catalog                               write Skill references (catalog-h5.md, catalog-android.md)
Phase 2  plan     --page P --platform X        reuse plan for a Figma page (gate: confirm before coding)
Phase 4  audit    --page P --platform X        static audit of the page source
         capture  --page P --platform X        render → screen.png + nodes.json
         diff     --page P --platform X        pixel + node diff vs Figma
         verify   --page P --platform X        audit + capture + diff + report (exit 3 if below threshold)
存量     match    --figma <export> <project>    new Figma page → which existing components to reuse and how to configure them
                  [--target xml|compose]       (.md / .json / .html overlay / component-map draft; no mapping needed)
         reuse    [dir] [--platform h5|android] brownfield health check of any Vue/React or Android project (no Figma needed)
                  [--focus files | --changed ref] page / PR gate: reinvented, missed reuse, bypassed wrappers (exit 3 = fail)
         compare  --url <page> --figma <export> any running H5 page vs the Figma frame: pixel diff + text coverage
dev      fixture  <design.html> -o <dir>       Figma-format export from an annotated design HTML

Global: run from the repo root (figmafe.yaml) or pass --root <dir>.
`)
}
