package core

import (
	"github.com/go-rod/rod/lib/launcher/flags"
	"os"
	"sync"
	"sync/atomic"

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

// poolSize is a counter for the number of browsers in the pool.
var poolSize int32

// maxPoolSize is the maximum size of the pool.
var maxPoolSize int32 = 10

func init() {
	browserPool = rod.NewBrowserPool(int(maxPoolSize))
}

// SetBrowserPool replaces the global browser pool.
func SetBrowserPool(pool rod.Pool[rod.Browser]) {
	if pool == nil {
		return
	}
	browserPool = pool
	atomic.StoreInt32(&poolSize, 0)
	maxPoolSize = int32(cap(pool))
}

// SetBrowserPoolSize resets the global pool size.
func SetBrowserPoolSize(size int) {
	if size <= 0 {
		return
	}
	browserPool = rod.NewBrowserPool(size)
	atomic.StoreInt32(&poolSize, 0)
	maxPoolSize = int32(size)
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
	l := newLauncherWithOptions(opts)
	launch := l.Bin(path).MustLaunch()

	var browser *rod.Browser
	if opts.UseGlobalPool {
		elem, err := browserPool.Get(func() (*rod.Browser, error) {
			connect := rod.New().ControlURL(launch).MustConnect()
			return connect, nil
		})
		if err != nil {
			browser = rod.New().MustIncognito().ControlURL(launch).MustConnect()
		} else {
			atomic.AddInt32(&poolSize, -1)
			browser = elem
		}
	} else {
		browser = rod.New().MustIncognito().ControlURL(launch).MustConnect()
	}

	return &Contain{
		Path:          path,
		UserDataDir:   opts.UserDataDir,
		Headless:      opts.Headless,
		LakeLess:      opts.Leakless,
		URL:           launch,
		useGlobalPool: opts.UseGlobalPool,
		Launcher:      l,
		Browser:       browser,
		Options:       cloneContainOptions(opts),
	}
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
	l := newLauncherWithOptions(opts)
	launch := l.Bin(path).MustLaunch()

	var browser *rod.Browser
	if pool != nil {
		elem, err := pool.Get(func() (*rod.Browser, error) {
			connect := rod.New().ControlURL(launch).MustConnect()
			return connect, nil
		})
		if err != nil {
			browser = rod.New().MustIncognito().ControlURL(launch).MustConnect()
		} else {
			browser = elem
		}
	} else {
		browser = rod.New().MustIncognito().ControlURL(launch).MustConnect()
	}

	return &Contain{
		Path:          path,
		UserDataDir:   opts.UserDataDir,
		Headless:      opts.Headless,
		LakeLess:      opts.Leakless,
		URL:           launch,
		useGlobalPool: false,
		Launcher:      l,
		Browser:       browser,
		Options:       cloneContainOptions(opts),
		CustomPool:    pool,
	}
}

// NewContainRemote connects to a remote Chrome instance.
func NewContainRemote(remoteURL string) *Contain {
	return NewContainRemoteWithPool(remoteURL, false)
}

// NewContainRemoteWithPool connects to remote Chrome with optional global pool usage.
func NewContainRemoteWithPool(remoteURL string, useGlobalPool bool) *Contain {
	var browser *rod.Browser
	if useGlobalPool {
		elem, err := browserPool.Get(func() (*rod.Browser, error) {
			connect := rod.New().ControlURL(remoteURL).MustConnect()
			return connect, nil
		})
		if err != nil {
			browser = rod.New().ControlURL(remoteURL).MustConnect()
		} else {
			atomic.AddInt32(&poolSize, -1)
			browser = elem
		}
	} else {
		browser = rod.New().ControlURL(remoteURL).MustConnect()
	}
	return &Contain{
		URL:           remoteURL,
		RemoteURL:     remoteURL,
		IsRemote:      true,
		useGlobalPool: useGlobalPool,
		Browser:       browser,
	}
}

// NewContainRemoteWithCustomPool connects to remote Chrome using a custom pool.
func NewContainRemoteWithCustomPool(remoteURL string, pool rod.Pool[rod.Browser]) *Contain {
	var browser *rod.Browser
	if pool != nil {
		elem, err := pool.Get(func() (*rod.Browser, error) {
			connect := rod.New().ControlURL(remoteURL).MustConnect()
			return connect, nil
		})
		if err != nil {
			browser = rod.New().ControlURL(remoteURL).MustConnect()
		} else {
			browser = elem
		}
	} else {
		browser = rod.New().ControlURL(remoteURL).MustConnect()
	}
	return &Contain{
		URL:        remoteURL,
		RemoteURL:  remoteURL,
		IsRemote:   true,
		Browser:    browser,
		CustomPool: pool,
	}
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

func createBrowserForPool(path, userDataDir string, headless bool, leakless bool) {
	currentSize := atomic.LoadInt32(&poolSize)
	if currentSize >= maxPoolSize {
		return
	}

	newSize := atomic.AddInt32(&poolSize, 1)
	if newSize > maxPoolSize {
		atomic.AddInt32(&poolSize, -1)
		return
	}

	var l *launcher.Launcher
	if userDataDir != "" {
		l = launcher.NewUserMode()
		l.Set("user-data-dir", userDataDir)
	} else {
		l = launcher.New()
	}
	launch := l.NoSandbox(true).Headless(headless).Leakless(leakless).
		Set("disable-setuid-sandbox").
		Set("disable-dev-shm-usage").
		Set("disable-accelerated-2d-canvas").
		Set("no-first-run").
		Set("no-zygote").
		Set("disable-gpu").
		Set("disable-web-security").
		Set("disable-infobars").
		Set("user-data-dir", os.TempDir()).
		Bin(path).
		MustLaunch()
	browser := rod.New().ControlURL(launch).MustConnect()
	browserPool.Put(browser)
}

func createBrowserForPoolWithOptions(path string, opts *ContainOptions) {
	if opts == nil {
		opts = DefaultContainOptions().WithGlobalPool(true)
	}
	currentSize := atomic.LoadInt32(&poolSize)
	if currentSize >= maxPoolSize {
		return
	}
	newSize := atomic.AddInt32(&poolSize, 1)
	if newSize > maxPoolSize {
		atomic.AddInt32(&poolSize, -1)
		return
	}
	launch := newLauncherWithOptions(opts).Bin(path).MustLaunch()
	browser := rod.New().ControlURL(launch).MustConnect()
	browserPool.Put(browser)
}

// Close closes the browser instance.
func (c *Contain) Close() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}

		path := c.Path
		headless := c.Headless
		leakless := c.LakeLess
		useGlobalPool := c.useGlobalPool
		opts := cloneContainOptions(c.Options)

		if c.Browser != nil {
			pages, err := c.Browser.Pages()
			if err == nil {
				for _, page := range pages {
					page.Close()
				}
			}
		}

		if c.CustomPool != nil {
			if c.Browser != nil {
				c.CustomPool.Put(c.Browser)
			}
			return
		}

		if c.Launcher != nil {
			c.Launcher.Cleanup()
			c.Launcher.Kill()
		}

		if c.IsRemote {
			if useGlobalPool && c.Browser != nil {
				browserPool.Put(c.Browser)
			}
			return
		}

		if useGlobalPool {
			go func() {
				if opts != nil {
					opts.WithGlobalPool(true)
					createBrowserForPoolWithOptions(path, opts)
					return
				}
				createBrowserForPool(path, c.UserDataDir, headless, leakless)
			}()
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
		useGlobalPool := c.useGlobalPool
		leakless := c.LakeLess
		opts := cloneContainOptions(c.Options)
		c.Close()

		if c.CustomPool != nil {
			if c.IsRemote && c.RemoteURL != "" {
				elem, err := c.CustomPool.Get(func() (*rod.Browser, error) {
					connect := rod.New().ControlURL(c.RemoteURL).MustConnect()
					return connect, nil
				})
				if err != nil {
					c.Browser = rod.New().ControlURL(c.RemoteURL).MustConnect()
				} else {
					c.Browser = elem
				}
				return
			}

			var newC *Contain
			if opts != nil {
				opts.WithGlobalPool(false)
				opts.WithHeadless(headless)
				opts.WithLeakless(leakless)
				newC = NewContainWithOptionsAndPool(path, opts, c.CustomPool)
			} else {
				newC = NewContainWithCustomPool(path, c.UserDataDir, headless, leakless, c.CustomPool)
			}
			*c = *newC
			return
		}

		if c.IsRemote && c.RemoteURL != "" {
			if useGlobalPool {
				elem, err := browserPool.Get(func() (*rod.Browser, error) {
					connect := rod.New().ControlURL(c.RemoteURL).MustConnect()
					return connect, nil
				})
				if err != nil {
					c.Browser = rod.New().ControlURL(c.RemoteURL).MustConnect()
				} else {
					atomic.AddInt32(&poolSize, -1)
					c.Browser = elem
				}
			} else {
				c.Browser = rod.New().ControlURL(c.RemoteURL).MustConnect()
			}
			return
		}

		var newC *Contain
		if opts != nil {
			opts.WithGlobalPool(useGlobalPool)
			opts.WithHeadless(headless)
			opts.WithLeakless(leakless)
			newC = NewContainWithOptions(path, opts)
		} else if useGlobalPool {
			newC = NewContain(path, c.UserDataDir, headless, leakless)
		} else {
			newC = NewContainWithoutPool(path, c.UserDataDir, headless, leakless)
		}
		*c = *newC
	})
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
