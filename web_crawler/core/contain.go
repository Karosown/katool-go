package core

import (
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/karosown/katool-go/lock"
)

// Contain 浏览器容器结构体
// Contain represents a browser container structure
type Contain struct {
	Path               string // 浏览器可执行文件路径 / Browser executable path
	Headless           bool   // 是否无头模式 / Whether headless mode
	URL                string // 启动器URL / Launcher URL
	*launcher.Launcher        // 启动器实例 / Launcher instance
	*rod.Browser              // 浏览器实例 / Browser instance
}

// WebReaderSysLock 全局读写锁
// WebReaderSysLock is a global read-write lock
var WebReaderSysLock *sync.RWMutex = &sync.RWMutex{}

// browserPool 全局浏览器池
// browserPool is a global browser pool
var browserPool rod.Pool[rod.Browser]

// init 初始化浏览器池
// init initializes the browser pool
func init() {
	browserPool = rod.NewBrowserPool(10) // 设置池大小为10
}

// NewCotain 创建新的Contain实例
// NewCotain creates a new Contain instance
func NewCotain(path string, headless bool) *Contain {
	l := launcher.NewUserMode()
	launch := l.NoSandbox(true).Headless(headless).Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-setuid-sandbox").
		Set("no-sandbox").
		Set("disable-web-security").
		Set("disable-infobars").Bin(path).MustLaunch()

	// 从浏览器池中获取浏览器实例
	elem, err := browserPool.Get(func() (*rod.Browser, error) {
		connect := rod.New().ControlURL(launch).MustConnect()
		return connect, nil
	})
	if err != nil {
		elem = rod.New().ControlURL(launch).MustConnect()
	}
	browser := elem

	return &Contain{
		Path:     path,
		Headless: headless,
		URL:      launch,
		Launcher: l,
		Browser:  browser,
	}
}

// GetContainer 获取浏览器实例
// GetContainer gets the browser instance
func (c *Contain) GetContainer() *rod.Browser {
	return c.Browser
}

// Close 关闭浏览器实例并归还到池中
// Close closes the browser instance and returns it to the pool
func (c *Contain) Close() {
	lock.Synchronized(WebReaderSysLock, func() {
		if c == nil {
			return
		}
		pages, err := c.Browser.Pages()
		if err == nil {
			for _, page := range pages {
				page.Close()
			}
		}
		// 将浏览器归还到池中
		c.Launcher.Cleanup()
		c.Launcher.Kill()
		if len(browserPool) < 3 {
			for len(browserPool) < 10 {
				go func() {
					if len(browserPool) > 10 {
						return
					}
					l := launcher.NewUserMode()
					launch := l.NoSandbox(true).Headless(c.Headless).Set("disable-gpu").
						Set("disable-dev-shm-usage").
						Set("disable-setuid-sandbox").
						Set("no-sandbox").
						Set("disable-web-security").
						Set("disable-infobars").Bin(c.Path).MustLaunch()
					browserPool.Put(rod.New().ControlURL(launch).MustConnect())
				}()
			}
		}
	})

}

// ReStart 重启浏览器
// ReStart restarts the browser
func (c *Contain) ReStart() {
	lock.Synchronized(WebReaderSysLock, func() {
		// 关闭当前实例
		c.Close()

		// 创建新实例并复制到当前对象
		newC := NewCotain(c.Path, c.Headless)
		*c = *newC
	})
}
