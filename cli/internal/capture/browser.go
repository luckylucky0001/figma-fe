package capture

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// ChromePath finds a Chromium binary: $FIGMAFE_CHROME, then the Playwright cache, then PATH lookup by chromedp.
func ChromePath() string {
	if p := os.Getenv("FIGMAFE_CHROME"); p != "" {
		return p
	}
	root := os.Getenv("PLAYWRIGHT_BROWSERS_PATH")
	if root == "" {
		root = "/opt/pw-browsers"
	}
	for _, pat := range []string{"chromium-*/chrome-linux/chrome", "chromium"} {
		m, _ := filepath.Glob(filepath.Join(root, pat))
		for _, p := range m {
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

// Browser is a headless Chromium session.
type Browser struct {
	Ctx    context.Context
	cancel func()
}

func NewBrowser(timeout time.Duration) (*Browser, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("font-render-hinting", "none"),
	)
	if p := ChromePath(); p != "" {
		opts = append(opts, chromedp.ExecPath(p))
	}
	actx, c1 := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, c2 := chromedp.NewContext(actx, chromedp.WithErrorf(func(string, ...any) {}))
	ctx, c3 := context.WithTimeout(ctx, timeout)
	if err := chromedp.Run(ctx); err != nil {
		c3()
		c2()
		c1()
		return nil, fmt.Errorf("start chromium: %w", err)
	}
	return &Browser{Ctx: ctx, cancel: func() { c3(); c2(); c1() }}, nil
}

func (b *Browser) Close() { b.cancel() }

// Open navigates to url with a mobile viewport of w×h CSS px at the given scale and waits for fonts.
func (b *Browser) Open(url string, w, h int, scale float64) error {
	return chromedp.Run(b.Ctx,
		emulation.SetDeviceMetricsOverride(int64(w), int64(h), scale, true),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(`document.fonts.ready.then(() => true)`, nil, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
		chromedp.Sleep(300*time.Millisecond),
	)
}

// Screenshot captures the viewport (w×h CSS px) at the device scale into a PNG file.
func (b *Browser) Screenshot(path string, w, h int, scale float64) error {
	var data string
	err := chromedp.Run(b.Ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		buf, err := page.CaptureScreenshot().
			WithFormat(page.CaptureScreenshotFormatPng).
			WithClip(&page.Viewport{X: 0, Y: 0, Width: float64(w), Height: float64(h), Scale: 1}).
			WithCaptureBeyondViewport(false).
			Do(ctx)
		if err != nil {
			return err
		}
		data = base64.StdEncoding.EncodeToString(buf)
		return nil
	}))
	if err != nil {
		return err
	}
	raw, _ := base64.StdEncoding.DecodeString(data)
	return os.WriteFile(path, raw, 0o644)
}

// Eval runs JS and unmarshals the JSON-serialisable result into out.
func (b *Browser) Eval(js string, out any) error {
	return chromedp.Run(b.Ctx, chromedp.Evaluate(js, out))
}
