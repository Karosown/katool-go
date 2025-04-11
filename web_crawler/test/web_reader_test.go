package test

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/web_crawler"
	"github.com/karosown/katool-go/web_crawler/core"
	"github.com/karosown/katool-go/web_crawler/render"
	"github.com/karosown/katool-go/web_crawler/rss"
)

func TestChorme(t *testing.T) {
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
	fmt.Printf("S  : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Printf("Data : %s\n", article.Content)
	fmt.Printf("TextContent : %s\n", article.SiteName)
	fmt.Printf("ImageiteName: %s\n", article.TextContent)
	fmt.Println()
}

func TestReadSourceCode(t *testing.T) {
	code := web_crawler.ReadSourceCode("https://openai.com/news/rss.xml", "", func(page *rod.Page) {
		page.MustStopLoading()
	})
	fmt.Println(code.String())
	r := &rss.RSS{}
	err := xml.Unmarshal([]byte(code.String()), r)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(r)
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
	u := "https://www.landiannews.com/feed"
	web_crawler.ReadRSS(u).Channel.Sites().Stream().Map(func(i rss.Item) any {
		path := web_crawler.ParsePath(u, i.Link)
		article := web_crawler.GetArticleWithChrome(path, nil, func(article web_crawler.Article) bool {
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
	web_crawler.WebChrome = core.NewCotain("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", true)
}
