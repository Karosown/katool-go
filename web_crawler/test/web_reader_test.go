package test

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/web_crawler"
	"github.com/karosown/katool-go/web_crawler/core"
	"github.com/karosown/katool-go/web_crawler/render"
	"github.com/karosown/katool-go/web_crawler/rss"
	"github.com/karosown/katool-go/words"
	"github.com/karosown/katool-go/words/split/jieba"
	"net/http"
	"strings"
	"testing"
	"unicode"
)

func TestChrome(t *testing.T) {
	container := launcher.MustResolveURL(web_crawler.ChromeRemoteURL)
	browser := rod.New().ControlURL(container).MustConnect().MustPage("https://openai.com/news/rss.xml")
	fmt.Println(
		browser.MustEval("() => {" +
			"	return document.documentElement.outerHTML" +
			"}"),
	)
}
func TestWebReader(t *testing.T) {
	u := "https://juejin.cn/post/7357922389417017353"
	article := web_crawler.GetArticleWithURL(u, func(r *http.Request) {
		r.Header = http.Header{
			"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""},
		}
	})
	if article.IsErr() {
		fmt.Println(article.Error())
	}
	fmt.Printf("URL     : %s\n", u)
	fmt.Printf("Title   : %s\n", article.Title)
	fmt.Printf("Author  : %s\n", article.Byline)
	fmt.Printf("Length  : %d\n", article.Length)
	fmt.Printf("Excerpt : %s\n", article.Excerpt)
	fmt.Printf("Image  : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Printf("Data : %s\n", article.Content)
	fmt.Printf("TextContent : %s\n", article.SiteName)
	fmt.Printf("ImageiteName: %s\n", article.TextContent)
	fmt.Println()
}

func TestReadChrome(t *testing.T) {
	for {
		u := "https://openai.com/index/openai-and-the-csu-system/"
		article := web_crawler.GetArticleWithChrome(u, func(page *rod.Page) {
			page.MustWaitLoad()
		}, nil)
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", u)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Data : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
	}
}
func TestReadClaude(t *testing.T) {
	//for {
	u := "https://www.anthropic.com/news"
	web_crawler.ReadArray(u, `()=>{
			// 使用 querySelectorAll 选择所有匹配的 <a> 标签  
			let linkElements = document.querySelectorAll('a.PostCard_post-card__z_Sqq.PostList_post-card__1g0fm');  
			
			// 将 NodeList 转换为数组并获取所有 href  
			let hrefValues = Array.from(linkElements).map(link => link.getArray('href'));  
			return hrefValues
		}`, func(page *rod.Page) {
		page.MustWaitLoad()
	}).Stream().Map(func(i web_crawler.WebReaderString) any {
		path := web_crawler.ParsePath(u, i.String())
		article := web_crawler.GetArticleWithURL(path, nil)
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", path)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Data : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
		return article
	}).ToList()

	//}
}
func TestReadOpenAI(t *testing.T) {
	//for {
	u := "https://openai.com/news/rss.xml"
	web_crawler.ReadRSS(u).Channel.Sites().Stream().Map(func(i rss.Item) any {
		path := web_crawler.ParsePath(u, i.Link)
		article := web_crawler.GetArticleWithChrome(path, func(page *rod.Page) {
			page.MustWaitLoad()
		}, func(article web_crawler.Article) bool {
			return article.Title == ""
		})
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", path)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Data : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
		fmt.Printf("PubTime:%s\n", article.PublishedTime)
		return article
	}).ToList()

	//}
}

func TestSourceRead(t *testing.T) {
	code := web_crawler.ReadSourceCode("https://www.cnbeta.com.tw/articles/tech/1483576.htm", "", render.Render)
	println(web_crawler.StripHTMLTags(web_crawler.DiyConvertHtmlToArticle(code.String())))
}
func init() {
	web_crawler.ChromeRemoteURL = "127.0.0.1:9222"
	web_crawler.WebChrome = core.NewCotain("C:\\Users\\Administrator\\AppData\\Local\\Google\\Chrome\\Bin\\Google Chrome.exe", true)
}

func TestWordAnalyzed(t *testing.T) {
	article := web_crawler.GetArticleWithChrome("https://www.36kr.com/p/3589639514423302", func(page *rod.Page) {
		page.MustWaitLoad()
	}, func(article web_crawler.Article) bool {
		return article.IsErr()
	})
	jieba.New().
		Cut(article.TextContent).
		Frequency().
		ToStream().
		Filter(func(i stream.Entry[string, int64]) bool {
			// 如果是纯英文
			if words.OnlyLanguage(i.Key, unicode.Han) {
				len := len([]rune(i.Key))
				return len > 3 && len < 6
			}
			return len(i.Key) > 4 && unicode.IsUpper(rune(i.Key[0])) && strings.ToLower(i.Key[1:]) == i.Key[1:]
		}).
		Filter(func(i stream.Entry[string, int64]) bool {
			return i.Value > 1
		}).
		Sort(func(a, b stream.Entry[string, int64]) bool {
			return a.Value > b.Value
		}).
		//Sub(0, 10).
		ForEach(func(item stream.Entry[string, int64]) {
			fmt.Println(item)
		})
}

func TestReadArr(t *testing.T) {
	web_crawler.ReadArray("https://36kr.com/", `()=> {
    // 1. 获取今天日期 (格式: 2025-12-10)
    const today = new Date().toISOString().split('T')[0];
    const urls = new Set();

    while (true) {
        let foundOld = false;
        const items = document.querySelectorAll('.kr-home-flow-item');
        
        // 2. 遍历当前页面所有文章
        for (const item of items) {
            const timeEl = item.querySelector('.kr-flow-bar-time');
            if (!timeEl) continue;

            const txt = timeEl.innerText.trim();
            const link = item.querySelector('.article-item-title')?.href;

            // 逻辑：包含"分钟/小时"或"今天的日期" -> 是今天文章
            if (txt.includes('分钟前') || txt.includes('小时前') || txt.includes(today)) {
                if (link) urls.add(link);
            } 
            // 逻辑：如果是日期格式且不包含今天 -> 是旧文章
            else if (/\d{4}-\d{2}-\d{2}/.test(txt)) {
                foundOld = true;
            }
        }

        // 3. 退出条件：找到旧文章 或 没有"查看更多"按钮
        const btn = document.querySelector('.kr-home-flow-see-more');
        if (foundOld || !btn || btn.offsetParent === null) break;

        // 4. 执行翻页并等待
        const countBefore = items.length;
        btn.click();
        
        // 轮询等待新元素出现（最多等 5 秒）
        for (let i = 0; i < 25; i++) {
            if (document.querySelectorAll('.kr-home-flow-item').length > countBefore) break;
        }
        
        // 如果点击后数量没变，认为到底了，退出
        if (document.querySelectorAll('.kr-home-flow-item').length === countBefore) break;
    }

    return Array.from(urls);
}
`, func(page *rod.Page) {
		page.MustWaitLoad()
	}).Stream().ForEach(func(item web_crawler.WebReaderString) {
		chrome := web_crawler.GetArticleWithChrome(item.String(), func(page *rod.Page) {
			page.MustWaitLoad()
		}, func(article web_crawler.Article) bool {
			return article.IsErr()
		})
		fmt.Println(chrome.Article.Title)
	})
}
