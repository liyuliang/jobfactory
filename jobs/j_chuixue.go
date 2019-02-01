package jobs

import (
	domParser "github.com/liyuliang/dom-parser"
	"github.com/liyuliang/utils/format"
	"github.com/liyuliang/models/protobuf"
	"time"
	"math"
	"log"
	"task-center/queue"
	"task-center/worker"
)

func addChuixueJobs() {
	categoryUrls := []string{}

	for c := 97; c < 123; c++ {
		ch := string(c)
		categoryUrls = append(categoryUrls, "http://www.chuixue.net/manhua/"+ch)
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

			total := dom.Find("#pager a b").Text()
			page := format.StrToInt(total) / 30

			page = int(math.Ceil(float64(page))) + 1

			bookUrls := []string{}

			bookUrls = append(bookUrls, categoryUrl)
			for i := 2; i <= page; i++ {
				bookUrls = append(bookUrls, categoryUrl+"/index_"+format.IntToStr(i)+".html")
			}
			println(page, len(bookUrls))

			for _, bookUrl := range bookUrls {

				println("checking book: ", bookUrl)
				html, err := HttpGet(bookUrl)
				if err != nil {
					println(err.Error())
					continue
				}
				html = gbkToUtf8(html)
				dom, err := domParser.InitDom(html)
				if err != nil {
					println(err.Error())
					continue
				}
				chapterUrls := []string{}
				for _, a := range dom.FindAll("#dmList ul li a.pic") {
					href, exist := a.Attr("href")
					if exist {
						chapterUrls = append(chapterUrls, "http://m.chuixue.net"+href)
					}
				}

				queueName := "parser_manhua_listing"
				models := []*queue.Model{}

				for _, url := range chapterUrls {
					site := getSite(url)
					if site == "" {
						log.Printf("Can not know site: %s", url)
						continue
					}

					m := protobuf.ParserManhuaPage{
						Site: site,
						Url:  url,
					}
					models = append(models, &queue.Model{
						Name:  queueName,
						Model: &m,
					})
				}

				worker.Pusher().New(models)
			}
		}
		time.Sleep(30 * time.Minute)
	}
}
