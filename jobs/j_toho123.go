package jobs

import (
	domParser "github.com/liyuliang/dom-parser"
	"fmt"
	"github.com/liyuliang/utils/format"
	"jobfactory/conf"
)

func addTohoJobs() {
	categoryUrls := []string{}

	for c := 97; c < 123; c++ {
		ch := string(c)

		for i := 1; i <= 20; i++ {
			url := fmt.Sprintf("https://m.tohomh123.com/f-1-----%s-updatetime--%d.html", ch, i)
			categoryUrls = append(categoryUrls, url)
		}
	}

	for _, categoryUrl := range categoryUrls {
		html, err := HttpGet(categoryUrl)
		if err != nil {
			println(err.Error())
			continue
		}

		dom, err := domParser.InitDom(html)
		if err != nil {
			println(err.Error())
		} else {

			bookUrls := []string{}
			for _, a := range dom.FindAll("ul.list li.am-thumbnail p.d-nowrap a") {
				href, exist := a.Attr("href")
				if exist {
					bookUrls = append(bookUrls, "https://m.tohomh123.com"+href)
				}
			}

			queueName := "parser_manhua_listing"
			for _, url := range bookUrls {

				api := fmt.Sprintf("%s?queue=%s&url=%s&delay=%d", conf.Remote().Get("api.job"), queueName, format.UrlEncode(url), randomSecond())
				html, err := HttpGet(api)
				if err != nil {
					println(err.Error())
				} else {
					println(html)
				}
			}
		}
	}
}
