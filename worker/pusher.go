package worker

import (
	"github.com/liyuliang/models/protobuf"
	"github.com/liyuliang/utils/format"
	"github.com/liyuliang/models/protobuf"
	"encoding/base64"
	"net/url"
	"log"
	"sync"
)

type Model struct {
	Id    uint64         `json:"Id,omitempty"`
	Name  string         `json:"Name"`
	Model protobuf.Model `json:"Model"`
}


type pusher struct {
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


func idsToVals(ids []string) (vals url.Values) {
	vals = url.Values{}
	for _, id := range ids {
		vals.Add("ids", id)
	}
	return vals
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
