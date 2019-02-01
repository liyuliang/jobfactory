package worker

import (
	"github.com/liyuliang/models/protobuf"
	"encoding/base64"
	"net/url"
	"log"
	"sync"
	"jobfactory/conf"
	"time"
	"bytes"
	"net/http"
	"errors"
	"io/ioutil"
)

type pusher struct {
	level     int
	delay     int
	sync.RWMutex
	dataGroup map[string][]string
}

func (p *pusher) addModels(queueName string, models []string) {
	if len(models) != 0 {
		vals := modelsToVals(models)

		addApi := conf.Remote().Get("api.queue") + queueName + conf.Remote().Get("api.queue.add")
		_, err := HttpAuthPost(addApi, vals)

		if err != nil {
			log.Printf("New %s jobs failed: %s", queueName, err.Error())
			_, err = HttpAuthPost(addApi, vals)
			if err != nil {
				log.Printf("ReNew %s jobs failed: %s", queueName, err.Error())
			} else {
				log.Printf("ReNew %s jobs by ids done .%d", queueName, len(models))
			}
		} else {
			log.Printf("New %s jobs by ids done .%d", queueName, len(models))
		}
	}
}
func (p *pusher) SetDelay(delay int) *pusher {
	p.delay = delay
	return p
}

func (p *pusher) SetLevel(level int) *pusher {
	p.level = level
	return p
}

func (p *pusher) New(jobs []*Model) {
	if len(jobs) != 0 {

		p.dataGroup = make(map[string][]string)

		for _, job := range jobs {
			queueName := job.Name

			b, _ := protobuf.Marshal(job.Model)
			encoded := base64.StdEncoding.EncodeToString([]byte(b))

			p.Lock()
			models := p.dataGroup[queueName]
			models = append(models, encoded)
			p.dataGroup[queueName] = models
			p.Unlock()
		}

		for queueName, models := range p.dataGroup {
			p.addModels(queueName, models)
		}
	}
}

func modelsToVals(models []string) (vals url.Values) {
	vals = url.Values{}
	for _, model := range models {
		vals.Add("models", model)
	}
	return vals
}

func Pusher() *pusher {

	p := new(pusher)
	return p
}

func HttpAuthPost(uri string, v url.Values) (content string, err error) {

	account := conf.Remote().Get("api.account")
	password := conf.Remote().Get("api.password")

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
