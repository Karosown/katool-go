package main

import (
	"codeberg.org/readeck/go-readability"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/file/file_serialize/csv"
	"github.com/karosown/katool-go/util/pathutil"
	"github.com/karosown/katool-go/web_crawler"
	"github.com/karosown/katool-go/web_crawler/core"
	"net/http"
	"net/url"
)

func main() {
	waiter := func(page *rod.Page) {
		page.WaitRepaint()
		page = page.MustWaitLoad()
	}
	u := "https://36kr.com/"
	hu, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	host := hu.Host
	list := stream.Em[web_crawler.ArticleLink, readability.Article](web_crawler.AllLink(
		//"curl \"https://www.46.la/hot\" ^\n  -H \"authority: www.46.la\" ^\n  -H \"accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\" ^\n  -H \"accept-language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6\" ^\n  -H \"cache-control: max-age=0\" ^\n  -H \"referer: https://www.bing.com/\" ^\n  -H ^\"sec-ch-ua: ^\\^\"Chromium^\\^\";v=^\\^\"122^\\^\", ^\\^\"Not(A:Brand^\\^\";v=^\\^\"24^\\^\", ^\\^\"Microsoft Edge^\\^\";v=^\\^\"122^\\^\"^\" ^\n  -H \"sec-ch-ua-mobile: ?0\" ^\n  -H ^\"sec-ch-ua-platform: ^\\^\"Windows^\\^\"^\" ^\n  -H \"sec-fetch-dest: document\" ^\n  -H \"sec-fetch-mode: navigate\" ^\n  -H \"sec-fetch-site: cross-site\" ^\n  -H \"sec-fetch-user: ?1\" ^\n  -H \"upgrade-insecure-requests: 1\" ^\n  -H \"user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0\"",
		//"curl \"https://36kr.com/\" ^\n  -H \"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\" ^\n  -H \"Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6\" ^\n  -H \"Cache-Control: max-age=0\" ^\n  -H \"Connection: keep-alive\" ^\n  -H ^\"Cookie: _waftokenid=eyJ2Ijp7ImEiOiIzbGRJMDBLWkZ3V25kZnNBMDVVS2xCdWpBYmpTcVJMMTZ3dlF6ZHJxNDBJPSIsImIiOjE3NjYzMjIyNzAsImMiOiJreXBBZWkrODFDVU9KdkZGaG5ZL1NocmY4TVFPRUNwZzJ4WUUyQWtDNU1ZPSJ9LCJzIjoiNVIvV2o4dnF1L1lpTWVkeDNtQ3JpZFBuTmJCYnIyUW9Uc0k4ZUFrYWNYaz0ifQ; sajssdk_2015_cross_new_user=1; sensorsdata2015jssdkcross=^%^7B^%^22distinct_id^%^22^%^3A^%^2219b4103496253a-0859bb42a3dc8b-4c657b58-2073600-19b4103496313d0^%^22^%^2C^%^22^%^24device_id^%^22^%^3A^%^2219b4103496253a-0859bb42a3dc8b-4c657b58-2073600-19b4103496313d0^%^22^%^2C^%^22props^%^22^%^3A^%^7B^%^7D^%^7D; SERVERID=6754aaff36cb16c614a357bbc08228ea^|1766322272^|1766322271^\" ^\n  -H \"Referer: https://36kr.com/\" ^\n  -H \"Sec-Fetch-Dest: document\" ^\n  -H \"Sec-Fetch-Mode: navigate\" ^\n  -H \"Sec-Fetch-Site: same-origin\" ^\n  -H \"Sec-Fetch-User: ?1\" ^\n  -H \"Upgrade-Insecure-Requests: 1\" ^\n  -H \"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0\" ^\n  -H ^\"sec-ch-ua: ^\\^\"Chromium^\\^\";v=^\\^\"122^\\^\", ^\\^\"Not(A:Brand^\\^\";v=^\\^\"24^\\^\", ^\\^\"Microsoft Edge^\\^\";v=^\\^\"122^\\^\"^\" ^\n  -H \"sec-ch-ua-mobile: ?0\" ^\n  -H ^\"sec-ch-ua-platform: ^\\^\"Windows^\\^\"^\"",
		u,
		waiter, func(r *http.Request) {
			r.Header = http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""}}
		}).Stream().Filter(func(i web_crawler.ArticleLink) bool {
		if len([]rune(i.Title)) < 5 {
			return false
		}
		parse, err := url.Parse(i.Link)
		if err != nil {
			return false
		}
		if parse.Host != host {
			return false
		}
		return len(pathutil.NewWrapper(i.Link).BeforeLayer().Path) > 8
	}).ForEach(func(item web_crawler.ArticleLink) {
		fmt.Println(item.Title, item.Link)
	}).Parallel()).Map(func(item web_crawler.ArticleLink) readability.Article {
		article := web_crawler.GetArticleWithChrome(item.Link, waiter, func(article web_crawler.Article) bool {
			return article.Title == ""
		}).Article
		println(item.Title)
		return article
	}).ToList()
	csv.WritePath("a.csv", list)
}
func init() {
	web_crawler.ChromeRemoteURL = "127.0.0.1:9222"
	web_crawler.ClosePageLimit = 300
	web_crawler.WebChrome = core.NewContainRemote("http://localhost:3000")
}
