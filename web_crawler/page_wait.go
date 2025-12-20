package web_crawler

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// PageWaitOptions controls how to wait for a page to finish async loading.
type PageWaitOptions struct {
	Selector    string
	MinCount    int
	ReadyState  string
	Timeout     time.Duration
	IdleDelay   time.Duration
	Scroll      bool
	MaxScroll   int
	ScrollDelay time.Duration
}

// DefaultPageWaitOptions returns default async wait options.
func DefaultPageWaitOptions() PageWaitOptions {
	return PageWaitOptions{
		ReadyState:  "complete",
		Timeout:     30 * time.Second,
		IdleDelay:   500 * time.Millisecond,
		MaxScroll:   10,
		ScrollDelay: 500 * time.Millisecond,
	}
}

// WaitPageReady waits for ready state and optional selector conditions.
func WaitPageReady(page *rod.Page, options PageWaitOptions) error {
	if page == nil {
		return errors.New("page is nil")
	}
	opts := options
	if opts.ReadyState == "" {
		opts.ReadyState = "complete"
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.IdleDelay < 0 {
		opts.IdleDelay = 0
	}
	if opts.ScrollDelay <= 0 {
		opts.ScrollDelay = 500 * time.Millisecond
	}
	if opts.MaxScroll < 0 {
		opts.MaxScroll = 0
	}

	if opts.Scroll {
		if err := autoScroll(page, opts); err != nil {
			return err
		}
	}

	minCount := opts.MinCount
	if opts.Selector != "" && minCount <= 0 {
		minCount = 1
	}

	deadline := time.Now().Add(opts.Timeout)
	for time.Now().Before(deadline) {
		readyOk := true
		if opts.ReadyState != "" {
			state, err := evalString(page, `() => document.readyState`)
			if err != nil {
				return err
			}
			readyOk = state == opts.ReadyState
		}

		selectorOk := true
		if opts.Selector != "" {
			count, err := evalInt(page, `(s) => document.querySelectorAll(s).length`, opts.Selector)
			if err != nil {
				return err
			}
			selectorOk = count >= minCount
		}

		if readyOk && selectorOk {
			if opts.IdleDelay > 0 {
				time.Sleep(opts.IdleDelay)
			}
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("wait page ready timeout after %s", opts.Timeout)
}

// MustWaitPageReady panics on wait error.
func MustWaitPageReady(page *rod.Page, options PageWaitOptions) {
	if err := WaitPageReady(page, options); err != nil {
		panic(err)
	}
}

// RenderWithWait returns a renderFunc using WaitPageReady.
func RenderWithWait(options PageWaitOptions) func(*rod.Page) {
	return func(page *rod.Page) {
		MustWaitPageReady(page, options)
	}
}

func autoScroll(page *rod.Page, options PageWaitOptions) error {
	maxScroll := options.MaxScroll
	if maxScroll <= 0 {
		maxScroll = 10
	}
	delay := options.ScrollDelay
	if delay <= 0 {
		delay = 500 * time.Millisecond
	}

	prevHeight := 0
	for i := 0; i < maxScroll; i++ {
		height, err := evalInt(page, `() => document.documentElement.scrollHeight`)
		if err != nil {
			return err
		}
		if height <= prevHeight && i > 0 {
			break
		}
		prevHeight = height
		if _, err := page.Eval(`(h) => window.scrollTo(0, h)`, height); err != nil {
			return err
		}
		time.Sleep(delay)
	}
	return nil
}

func evalString(page *rod.Page, js string, args ...interface{}) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("eval string failed: %v", r)
		}
	}()
	val := page.MustEval(js, args...)
	return val.Str(), nil
}

func evalInt(page *rod.Page, js string, args ...interface{}) (res int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("eval int failed: %v", r)
		}
	}()
	val := page.MustEval(js, args...)
	return val.Int(), nil
}
