package jobs

import (
	domParser "github.com/liyuliang/dom-parser"
	"time"
	"fmt"
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
					bookUrls = append(bookUrls, "https://m.tohomh123.com" + href)
				}
			}
		}

		time.Sleep(30 * time.Minute)
	}
}
