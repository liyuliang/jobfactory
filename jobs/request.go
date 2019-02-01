package jobs

import (
	"github.com/liyuliang/utils/regex"
	"github.com/imroc/req"
	"time"
	"bytes"
	"net/http"
	"log"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
	"github.com/golang/text/transform"
	"github.com/golang/text/encoding/simplifiedchinese"
)

func gbkToUtf8(text string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(text)), simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(reader)
	return string(d)
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

func HttpGet(uri string) (string, error) {

	Url, err := url.Parse(uri)
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

func HttpAuthPost(uri string, v url.Values) (content string, err error) {

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
