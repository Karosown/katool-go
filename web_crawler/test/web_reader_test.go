package test

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
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
	"time"
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
	web_crawler.WebChrome = core.NewCotain("C:\\Program Files\\Google\\Chrome\\Application\\Chrome.exe", false)
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

func TestReadLink(t *testing.T) {
	web_crawler.ReadArticleLinksAuto("https://www.zhihu.com/", func(page *rod.Page) {
		page.MustEvalOnNewDocument(stealth.JS)
		page.MustWaitLoad()
		page.MustWaitRequestIdle()
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""})
		page.MustElement("body").MustHTML()
		page.WaitIdle(5 * time.Second)
	}, func(r *http.Request) {
		r.Header = http.Header{
			"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""},
		}
	}).Stream().ForEach(func(item web_crawler.ArticleLink) {
		fmt.Println(item)
	})
}
func TestGenaerte(t *testing.T) {
	snippet, err := web_crawler.GenerateRodPageSetupSnippetFromCURL("curl \"https://www.zhihu.com/hot\" ^\n  -H \"authority: www.zhihu.com\" ^\n  -H \"accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\" ^\n  -H \"accept-language: zh-CN,zh;q=0.9,en;q=0.8\" ^\n  -H \"cache-control: max-age=0\" ^\n  -H \"cookie: _zap=d495d78c-7631-4845-a130-5f004e24afb3; d_c0=iIbTlMGuQBuPTjTVO9CyIa2kngnD41bLNYk=^|1760977246; __snaker__id=nGDK5t0boNwgFCJS; q_c1=d316d753f0a74ef8ad2fdd2344912ef5^|1760977283000^|1760977283000; _xsrf=fGHZ2JNZuQi7pJcKkkjLOhd6rDcn1yFz; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1764082829,1764863435; z_c0=2^|1:0^|10:1765377656^|4:z_c0^|92:Mi4xUnR0OUN3QUFBQUNJaHRPVXdhNUFHeVlBQUFCZ0FsVk4wNDBTYWdEb3BBWEhYWDhZTGthUDVvTW1TTEZkU0gtTFRB^|c20bcb2ca37e3641dfb3e09cd80485e11887a804e754ba465ad53f3c1fcab271; __zse_ck=004_d/1en7XDkyGmw=o3zqJVzA10X2rDt=VM0iEJRiQSEvTiCEOUQWW9A6T9jnbGs5oF5zmDxWY9KgZTMD5qgnAzMuFrZenkfi3RFBdMkyUKSmPUXU2Ny5OtX7EitXGOGNLi-7tq72p4IIo1epXYqrZGQ672e92DUL873et2GE8zFYCeLf8A5N718ZrZ+OaTUnMuRus6htbWWYkUYrsmi7JYEe2BJTLm2pOvs/umQLSsxg4N9uk8KfoJiVoRBMLyhw0Nf; BEC=fc13dc7850b2e749d88c66e883fdd0e4; SESSIONID=qhx0ZZvBjgwxvGNTGrilNEiRHdHBCoKk3FEgHNRyxh0; JOID=UFEWB08mVJTm_pEKRbLIRlI7OUdfXDbr1qLAdgBKOOmZiNBPG6st1YnzkAdCbwM7WM30T0IKS7L3YBQ5izRy0gE=; osd=UVEQBk0nVJLn_JAKQ7PKR1I9OEVeXDDq1KPAcAFIOemfidJOG60s14jzlgZAbgM9Wc_1T0QLSbP3ZhU7ijR00wM=\" ^\n  -H \"sec-ch-ua: ^\\^\"Not_A Brand^\\^\";v=^\\^\"99^\\^\", ^\\^\"Google Chrome^\\^\";v=^\\^\"109^\\^\", ^\\^\"Chromium^\\^\";v=^\\^\"109^\\^\"\" ^\n  -H \"sec-ch-ua-mobile: ?0\" ^\n  -H \"sec-ch-ua-platform: ^\\^\"Windows^\\^\"\" ^\n  -H \"sec-fetch-dest: document\" ^\n  -H \"sec-fetch-mode: navigate\" ^\n  -H \"sec-fetch-site: none\" ^\n  -H \"sec-fetch-user: ?1\" ^\n  -H \"upgrade-insecure-requests: 1\" ^\n  -H \"user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36\" ^\n  --compressed")
	if err != nil {
		panic(err)
	}
	fmt.Println(snippet)

}
