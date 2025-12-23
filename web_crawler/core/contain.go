package core

import (
	"github.com/go-rod/rod/lib/launcher/flags"
	"os"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/karosown/katool-go/lock"
)

// Contain represents a browser container.
type Contain struct {
	Path          string
	UserDataDir   string
	Headless      bool
	URL           string
	useGlobalPool bool
	*launcher.Launcher
	*rod.Browser
	LakeLess   bool
	RemoteURL  string
	IsRemote   bool
	Options    *ContainOptions
	CustomPool rod.Pool[rod.Browser]
}

// ContainOptions provides customization for launching Chrome.
type ContainOptions struct {
	Headless       bool
	Leakless       bool
	UseGlobalPool  bool
	UseUserMode    bool
	UseDefaultArgs bool
	UserDataDir    string
	UserAgent      string
	Flags          []string
	Args           map[string]string
}

// DefaultContainOptions returns default options.
func DefaultContainOptions() *ContainOptions {
	return &ContainOptions{
		Headless:       false,
		Leakless:       false,
		UseGlobalPool:  true,
		UseUserMode:    true,
		UseDefaultArgs: true,
		UserDataDir:    os.TempDir(),
		Flags:          nil,
		Args:           map[string]string{},
	}
}

func (o *ContainOptions) WithHeadless(v bool) *ContainOptions {
	o.Headless = v
	return o
}

func (o *ContainOptions) WithLeakless(v bool) *ContainOptions {
	o.Leakless = v
	return o
}

func (o *ContainOptions) WithGlobalPool(v bool) *ContainOptions {
	o.UseGlobalPool = v
	return o
}

func (o *ContainOptions) WithUserMode(v bool) *ContainOptions {
	o.UseUserMode = v
	return o
}

func (o *ContainOptions) WithDefaultArgs(v bool) *ContainOptions {
	o.UseDefaultArgs = v
	return o
}

func (o *ContainOptions) WithUserDataDir(dir string) *ContainOptions {
	o.UserDataDir = dir
	return o
}

func (o *ContainOptions) WithUserAgent(ua string) *ContainOptions {
	o.UserAgent = ua
	return o
}

func (o *ContainOptions) AddFlag(flag string) *ContainOptions {
	if flag == "" {
		return o
	}
	o.Flags = append(o.Flags, flag)
	return o
}

func (o *ContainOptions) SetArg(key, val string) *ContainOptions {
	if key == "" {
		return o
	}
	if o.Args == nil {
		o.Args = map[string]string{}
	}
	o.Args[key] = val
	return o
}

// WebReaderSysLock is a global read-write lock.
var WebReaderSysLock *sync.RWMutex = &sync.RWMutex{}

// browserPool is a global browser pool.
var browserPool rod.Pool[rod.Browser]

const defaultPoolSize = 10

func init() {
	browserPool = rod.NewBrowserPool(defaultPoolSize)
}

// SetBrowserPool replaces the global browser pool.
func SetBrowserPool(pool rod.Pool[rod.Browser]) {
	if pool == nil {
		return
	}
	browserPool = pool
}

// SetBrowserPoolSize resets the global pool size.
func SetBrowserPoolSize(size int) {
	if size <= 0 {
		return
	}
	browserPool = rod.NewBrowserPool(size)
}

// NewCotain creates a new Contain instance using default options.
func NewCotain(path string, headlessAndLeakLess ...bool) *Contain {
	headless := false
	leakless := false
	if len(headlessAndLeakLess) > 0 {
		headless = headlessAndLeakLess[0]
	}
	if len(headlessAndLeakLess) > 1 {
		leakless = headlessAndLeakLess[1]
	}
	opts := DefaultContainOptions().
		WithHeadless(headless).
		WithLeakless(leakless).
		WithGlobalPool(true)
	return NewContainWithOptions(path, opts)
}

// NewContain creates a new Contain instance using the global pool.
func NewContain(path, userDataDir string, headless, leakless bool) *Contain {
	opts := DefaultContainOptions().
		WithHeadless(headless).
		WithLeakless(leakless).
		WithGlobalPool(true)
	if userDataDir != "" {
		opts.WithUserDataDir(userDataDir).WithUserMode(true)
	} else {
		opts.WithUserDataDir("").WithUserMode(false)
	}
	return NewContainWithOptions(path, opts)
}

// NewContainWithoutPool creates a new Contain instance without using the global pool.
func NewContainWithoutPool(path, userDataDir string, headless bool, leakless bool) *Contain {
	opts := DefaultContainOptions().
		WithHeadless(headless).
		WithLeakless(leakless).
		WithGlobalPool(false)
	if userDataDir != "" {
		opts.WithUserDataDir(userDataDir).WithUserMode(true)
	} else {
		opts.WithUserDataDir("").WithUserMode(false)
	}
	return NewContainWithOptions(path, opts)
}

// NewContainWithPool creates a new Contain instance with optional global pool usage.
func NewContainWithPool(path, userDataDir string, headless bool, useGlobalPool bool, leakless bool) *Contain {
	opts := DefaultContainOptions().
		WithHeadless(headless).
		WithLeakless(leakless).
		WithGlobalPool(useGlobalPool)
	if userDataDir != "" {
		opts.WithUserDataDir(userDataDir).WithUserMode(true)
	} else {
		opts.WithUserDataDir("").WithUserMode(false)
	}
	return NewContainWithOptions(path, opts)
}

// NewContainWithOptions creates a new Contain instance with options.
func NewContainWithOptions(path string, options *ContainOptions) *Contain {
	opts := cloneContainOptions(options)
	if opts == nil {
		opts = DefaultContainOptions()
	}
	pool := rod.Pool[rod.Browser](nil)
	if opts.UseGlobalPool {
		pool = browserPool
	}
	return newContainWithPool(path, opts, pool, opts.UseGlobalPool, nil)
}

// NewContainWithCustomPool creates a new Contain instance using a custom pool.
func NewContainWithCustomPool(path, userDataDir string, headless bool, leakless bool, pool rod.Pool[rod.Browser]) *Contain {
	opts := DefaultContainOptions().
		WithHeadless(headless).
		WithLeakless(leakless).
		WithGlobalPool(false)
	if userDataDir != "" {
		opts.WithUserDataDir(userDataDir).WithUserMode(true)
	} else {
		opts.WithUserDataDir("").WithUserMode(false)
	}
	return NewContainWithOptionsAndPool(path, opts, pool)
}

// NewContainWithOptionsAndPool creates a new Contain instance with options and a custom pool.
func NewContainWithOptionsAndPool(path string, options *ContainOptions, pool rod.Pool[rod.Browser]) *Contain {
	opts := cloneContainOptions(options)
	if opts == nil {
		opts = DefaultContainOptions()
	}
	return newContainWithPool(path, opts, pool, false, pool)
}

// NewContainRemote connects to a remote Chrome instance.
func NewContainRemote(remoteURL string) *Contain {
	return NewContainRemoteWithPool(remoteURL, false)
}

// NewContainRemoteWithPool connects to remote Chrome with optional global pool usage.
func NewContainRemoteWithPool(remoteURL string, useGlobalPool bool) *Contain {
	pool := rod.Pool[rod.Browser](nil)
	if useGlobalPool {
		pool = browserPool
	}
	return newRemoteContain(remoteURL, pool, useGlobalPool, nil)
}

// NewContainRemoteWithCustomPool connects to remote Chrome using a custom pool.
func NewContainRemoteWithCustomPool(remoteURL string, pool rod.Pool[rod.Browser]) *Contain {
	return newRemoteContain(remoteURL, pool, false, pool)
}

// GetContainer gets the browser instance.
func (c *Contain) GetContainer() *rod.Browser {
	return c.Browser
}

// PageWithStealth creates a page with stealth mode.
func (c *Contain) PageWithStealth(url string) *rod.Page {
	page := stealth.MustPage(c.Browser)
	if url != "" {
		page.MustNavigate(url)
	}
	return page
}

// Close closes the browser instance or releases it back to the pool.
func (c *Contain) Close() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}

		if c.Browser != nil {
			pages, err := c.Browser.Pages()
			if err == nil {
				for _, page := range pages {
					page.Close()
				}
			}
		}

		if c.CustomPool != nil || c.useGlobalPool {
			pool := c.CustomPool
			if pool == nil {
				pool = browserPool
			}
			if c.Browser != nil {
				_ = c.Browser.Close()
			}
			if c.Launcher != nil {
				c.Launcher.Cleanup()
				c.Launcher.Kill()
			}
			var newBrowser *rod.Browser
			if c.IsRemote && c.RemoteURL != "" {
				newBrowser = rod.New().ControlURL(c.RemoteURL).MustConnect()
			} else {
				opts := cloneContainOptions(c.Options)
				if opts == nil {
					opts = DefaultContainOptions()
				}
				opts.WithHeadless(c.Headless)
				opts.WithLeakless(c.LakeLess)
				if c.UserDataDir != "" {
					opts.WithUserDataDir(c.UserDataDir).WithUserMode(true)
				} else {
					opts.WithUserDataDir("").WithUserMode(false)
				}
				newBrowser, _ = createLocalBrowser(c.Path, opts)
			}
			pool.Put(newBrowser)
			c.Browser = nil
			return
		}

		if c.Launcher != nil {
			c.Launcher.Cleanup()
			c.Launcher.Kill()
		}
	})
}

// ReStart restarts the browser.
func (c *Contain) ReStart() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}

		path := c.Path
		headless := c.Headless
		leakless := c.LakeLess
		opts := cloneContainOptions(c.Options)
		remoteURL := c.RemoteURL
		isRemote := c.IsRemote

		c.Close()

		if c.CustomPool != nil {
			if isRemote && remoteURL != "" {
				newC := newRemoteContain(remoteURL, c.CustomPool, false, c.CustomPool)
				*c = *newC
				return
			}
			if opts != nil {
				opts.WithGlobalPool(false)
				opts.WithHeadless(headless)
				opts.WithLeakless(leakless)
				newC := NewContainWithOptionsAndPool(path, opts, c.CustomPool)
				*c = *newC
				return
			}
			newC := NewContainWithCustomPool(path, c.UserDataDir, headless, leakless, c.CustomPool)
			*c = *newC
			return
		}

		if c.useGlobalPool {
			if isRemote && remoteURL != "" {
				newC := newRemoteContain(remoteURL, browserPool, true, nil)
				*c = *newC
				return
			}
			if opts != nil {
				opts.WithGlobalPool(true)
				opts.WithHeadless(headless)
				opts.WithLeakless(leakless)
				newC := NewContainWithOptions(path, opts)
				*c = *newC
				return
			}
			newC := NewContain(path, c.UserDataDir, headless, leakless)
			*c = *newC
			return
		}

		if isRemote && remoteURL != "" {
			newC := newRemoteContain(remoteURL, nil, false, nil)
			*c = *newC
			return
		}

		if opts != nil {
			opts.WithGlobalPool(false)
			opts.WithHeadless(headless)
			opts.WithLeakless(leakless)
			newC := NewContainWithOptions(path, opts)
			*c = *newC
			return
		}

		newC := NewContainWithoutPool(path, c.UserDataDir, headless, leakless)
		*c = *newC
	})
}

func newContainWithPool(path string, opts *ContainOptions, pool rod.Pool[rod.Browser], useGlobalPool bool, customPool rod.Pool[rod.Browser]) *Contain {
	browser, launcherInst := getBrowserFromPool(pool, func() (*rod.Browser, *launcher.Launcher) {
		return createLocalBrowser(path, opts)
	})
	launchURL := ""
	if launcherInst != nil {
		launchURL = launcherInst.MustLaunch()
	}

	return &Contain{
		Path:          path,
		UserDataDir:   opts.UserDataDir,
		Headless:      opts.Headless,
		LakeLess:      opts.Leakless,
		URL:           launchURL,
		useGlobalPool: useGlobalPool,
		Launcher:      launcherInst,
		Browser:       browser,
		Options:       cloneContainOptions(opts),
		CustomPool:    customPool,
	}
}

func newRemoteContain(remoteURL string, pool rod.Pool[rod.Browser], useGlobalPool bool, customPool rod.Pool[rod.Browser]) *Contain {
	browser, _ := getBrowserFromPool(pool, func() (*rod.Browser, *launcher.Launcher) {
		return rod.New().ControlURL(remoteURL).MustConnect(), nil
	})

	return &Contain{
		URL:           remoteURL,
		RemoteURL:     remoteURL,
		IsRemote:      true,
		useGlobalPool: useGlobalPool,
		Browser:       browser,
		CustomPool:    customPool,
	}
}

func getBrowserFromPool(pool rod.Pool[rod.Browser], create func() (*rod.Browser, *launcher.Launcher)) (*rod.Browser, *launcher.Launcher) {
	if pool == nil {
		return create()
	}
	var createdLauncher *launcher.Launcher
	elem, err := pool.Get(func() (*rod.Browser, error) {
		browser, l := create()
		createdLauncher = l
		return browser, nil
	})
	if err != nil {
		return create()
	}
	return elem, createdLauncher
}

func createLocalBrowser(path string, opts *ContainOptions) (*rod.Browser, *launcher.Launcher) {
	l := newLauncherWithOptions(opts)
	launch := l.Bin(path).MustLaunch()
	browser := rod.New().MustIncognito().ControlURL(launch).MustConnect()
	return browser, l
}

func newLauncherWithOptions(opts *ContainOptions) *launcher.Launcher {
	if opts == nil {
		opts = DefaultContainOptions()
	}
	var l *launcher.Launcher
	if opts.UseUserMode {
		l = launcher.NewUserMode()
	} else {
		l = launcher.New()
	}
	l = l.NoSandbox(true).
		Leakless(opts.Leakless).
		Headless(opts.Headless)

	if opts.UseDefaultArgs {
		l = l.Set("disable-setuid-sandbox").
			Set("disable-dev-shm-usage").
			Set("disable-accelerated-2d-canvas").
			Set("no-first-run").
			Set("no-zygote").
			Set("disable-gpu").
			Set("disable-web-security").
			Set("disable-infobars")
	}
	if opts.UserDataDir != "" {
		l = l.Set("user-data-dir", opts.UserDataDir)
	}
	if opts.UserAgent != "" {
		l = l.Set("user-agent", opts.UserAgent)
	}
	for _, flag := range opts.Flags {
		if flag == "" {
			continue
		}
		l = l.Set(flags.Flag(flag))
	}
	for k, v := range opts.Args {
		if k == "" {
			continue
		}
		l = l.Set(flags.Flag(k), v)
	}
	return l
}

func cloneContainOptions(src *ContainOptions) *ContainOptions {
	if src == nil {
		return nil
	}
	dst := *src
	if src.Flags != nil {
		dst.Flags = append([]string{}, src.Flags...)
	}
	if src.Args != nil {
		dst.Args = make(map[string]string, len(src.Args))
		for k, v := range src.Args {
			dst.Args[k] = v
		}
	}
	return &dst
}
