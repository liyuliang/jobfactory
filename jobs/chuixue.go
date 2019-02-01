package jobs

import (
	UrlPkg "net/url"
	domParser "github.com/liyuliang/dom-parser"
	"github.com/liyuliang/utils/regex"
	"github.com/imroc/req"
	"github.com/liyuliang/utils/format"
	"github.com/golang/text/transform"
	"github.com/golang/text/encoding/simplifiedchinese"
	"github.com/liyuliang/models/protobuf"
	"io/ioutil"
	"time"
	"math"
	"bytes"
	"strings"
	"log"
	"net/http"
	"errors"
	"jobfactory/worker"
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
				models := []*worker.Model{}

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
					models = append(models, &worker.Model{
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

func getSite(uri string) (site string) {
	if strings.Contains(uri, "chuixue") {
		site = "chuixue"
	}
	if strings.Contains(uri, "dajiaochong") {
		site = "dajiaochong"
	}
	if strings.Contains(uri, "dmzj") {
		site = "dmzj"
	}
	return site
}

func gbkToUtf8(text string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(text)), simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(reader)
	return string(d)
}

func HttpGet(uri string) (string, error) {

	Url, err := UrlPkg.Parse(uri)
	if err != nil {
	}

	referer := regex.Get(Url.Host, `([^\.]+\.[^\.]+$)`)
	referer = Url.Scheme + "://www." + referer

	r := req.New()

	r.EnableInsecureTLS(true)
	r.SetTimeout(10 * time.Second)

	header := "Mozilla/5.0 (iPhone; CPU iPhone OS 5_1_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9B206 Safari/7534.48.3"
	resp, err := r.Get(uri, req.Header{
		"User-Agent": header,
		"Referer":    referer,
	})

	if err != nil {
		return "", err
	}

	if resp.Response().StatusCode != 200 {
		// Try again
		resp, err = r.Get(uri, req.Header{
			"User-Agent": header,
			"Referer":    referer,
		})
		if err != nil {
			return "", err
		}

		if resp.Response().StatusCode != 200 {
			return "", err
		}
	}

	return resp.String(), nil
}

func HttpAuthPost(uri string, v UrlPkg.Values) (content string, err error) {

	account := "liang"
	password := "L!ang*#06#"

	p := bytes.NewBufferString(v.Encode())

	requeset, err := http.NewRequest("POST", uri, p)
	requeset.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requeset.SetBasicAuth(account, password)

	client := &http.Client{
		Timeout: time.Duration(60 * time.Second),
	}
	resp, err := client.Do(requeset)
	if err != nil {
		log.Println(err.Error())
		return "", err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return "", errors.New("Http response code is not 200. ")
		} else {
			bodyText, err := ioutil.ReadAll(resp.Body)
			return string(bodyText), err
		}
	}
}
