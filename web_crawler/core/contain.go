package core

import (
	"sync"
	"sync/atomic"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/karosown/katool-go/lock"
)

// Contain 浏览器容器结构体
// Contain represents a browser container structure
type Contain struct {
	Path               string // 浏览器可执行文件路径 / Browser executable path
	Headless           bool   // 是否无头模式 / Whether headless mode
	URL                string // 启动器URL / Launcher URL
	useGlobalPool      bool   // 是否使用全局池 / Whether to use global pool
	*launcher.Launcher        // 启动器实例 / Launcher instance
	*rod.Browser              // 浏览器实例 / Browser instance
}

// WebReaderSysLock 全局读写锁
// WebReaderSysLock is a global read-write lock
var WebReaderSysLock *sync.RWMutex = &sync.RWMutex{}

// browserPool 全局浏览器池
// browserPool is a global browser pool
var browserPool rod.Pool[rod.Browser]

// poolSize 池中浏览器数量计数器（用于跟踪池的大小）
// poolSize is a counter for the number of browsers in the pool (used to track pool size)
var poolSize int32

// maxPoolSize 池的最大大小
// maxPoolSize is the maximum size of the pool
const maxPoolSize = 10

// init 初始化浏览器池
// init initializes the browser pool
func init() {
	browserPool = rod.NewBrowserPool(maxPoolSize) // 设置池大小为10
}

// NewCotain 创建新的Contain实例（使用全局池）
// NewCotain creates a new Contain instance (using global pool)
func NewCotain(path string, headless bool) *Contain {
	return NewContainWithPool(path, headless, true)
}

// NewContainWithoutPool 创建新的Contain实例（不使用全局池，独立管理浏览器）
// NewContainWithoutPool creates a new Contain instance (without global pool, manages browser independently)
func NewContainWithoutPool(path string, headless bool) *Contain {
	return NewContainWithPool(path, headless, false)
}

// NewContainWithPool 创建新的Contain实例（可选择是否使用全局池）
// NewContainWithPool creates a new Contain instance (with option to use global pool)
func NewContainWithPool(path string, headless bool, useGlobalPool bool) *Contain {
	l := launcher.NewUserMode()
	launch := l.NoSandbox(true).Headless(headless).Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-setuid-sandbox").
		Set("no-sandbox").
		Set("disable-web-security").
		Set("disable-infobars").Bin(path).MustLaunch()

	var browser *rod.Browser

	if useGlobalPool {
		// 从全局浏览器池中获取浏览器实例
		elem, err := browserPool.Get(func() (*rod.Browser, error) {
			// 池为空时创建新浏览器（不减少计数器，因为这是新创建的）
			connect := rod.New().ControlURL(launch).MustConnect()
			return connect, nil
		})
		if err != nil {
			// 如果从池中获取失败，直接创建新浏览器
			browser = rod.New().ControlURL(launch).MustConnect()
		} else {
			// 成功从池中获取，减少计数器
			atomic.AddInt32(&poolSize, -1)
			browser = elem
		}
	} else {
		// 不使用全局池，直接创建新浏览器
		browser = rod.New().ControlURL(launch).MustConnect()
	}

	return &Contain{
		Path:          path,
		Headless:      headless,
		URL:           launch,
		useGlobalPool: useGlobalPool,
		Launcher:      l,
		Browser:       browser,
	}
}

// GetContainer 获取浏览器实例
// GetContainer gets the browser instance
func (c *Contain) GetContainer() *rod.Browser {
	return c.Browser
}

// PageWithStealth 创建一个带有Stealth模式的页面
// PageWithStealth creates a page with Stealth mode
func (c *Contain) PageWithStealth(url string) *rod.Page {
	page := stealth.MustPage(c.Browser)
	if url != "" {
		page.MustNavigate(url)
	}
	return page
}

// createBrowserForPool 创建新浏览器并放入池中（仅在池未满时）
// createBrowserForPool creates a new browser and puts it into the pool (only if pool is not full)
func createBrowserForPool(path string, headless bool) {
	// 检查池是否需要补充
	currentSize := atomic.LoadInt32(&poolSize)
	if currentSize >= maxPoolSize {
		// 池已满，不需要创建新浏览器
		return
	}

	// 尝试增加计数器（乐观锁）
	newSize := atomic.AddInt32(&poolSize, 1)
	if newSize > maxPoolSize {
		// 如果超过最大大小，回退计数器并返回
		atomic.AddInt32(&poolSize, -1)
		return
	}

	// 创建新浏览器
	l := launcher.NewUserMode()
	launch := l.NoSandbox(true).Headless(headless).Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-setuid-sandbox").
		Set("no-sandbox").
		Set("disable-web-security").
		Set("disable-infobars").Bin(path).MustLaunch()
	browser := rod.New().ControlURL(launch).MustConnect()

	// 放入池中（rod.Pool 会自动处理容量限制）
	browserPool.Put(browser)
}

// Close 关闭浏览器实例
// 如果使用全局池：kill后创建新浏览器放入池中
// 如果不使用全局池：直接kill浏览器，不创建新的
// Close closes the browser instance
// If using global pool: kills it, then creates a new browser and puts it into the pool
// If not using global pool: directly kills the browser without creating a new one
func (c *Contain) Close() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}

		// 保存配置信息
		path := c.Path
		headless := c.Headless
		useGlobalPool := c.useGlobalPool

		// 关闭所有页面
		if c.Browser != nil {
			pages, err := c.Browser.Pages()
			if err == nil {
				for _, page := range pages {
					page.Close()
				}
			}
		}

		// Kill 当前浏览器和启动器
		if c.Launcher != nil {
			c.Launcher.Cleanup()
			c.Launcher.Kill()
		}

		// 如果使用全局池，创建新浏览器并放入池中（异步执行，避免阻塞）
		// 如果不使用全局池，直接结束，不创建新浏览器
		if useGlobalPool {
			go func() {
				createBrowserForPool(path, headless)
			}()
		}
	})
}

// ReStart 重启浏览器
// ReStart restarts the browser
func (c *Contain) ReStart() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}

		// 保存配置信息
		path := c.Path
		headless := c.Headless
		useGlobalPool := c.useGlobalPool

		// 关闭当前实例
		c.Close()

		// 根据配置创建新实例并复制到当前对象
		var newC *Contain
		if useGlobalPool {
			newC = NewCotain(path, headless)
		} else {
			newC = NewContainWithoutPool(path, headless)
		}
		*c = *newC
	})
}
