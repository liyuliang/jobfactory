package controllers

import (
	"github.com/astaxie/beego"
	"github.com/liyuliang/models/protobuf"
	"jobfactory/worker"
	"github.com/liyuliang/utils/format"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	url := c.GetString("url")
	queueName := c.GetString("queue")
	level, _ := c.GetInt("level")
	delay, _ := c.GetInt("delay")

	if url == "" {
		c.Abort("404")
		return
	}
	url, _ = format.UrlDecode(url)

	models := []*worker.Model{}

	m := protobuf.ParserManhuaPage{
		Url: url,
	}

	models = append(models, &worker.Model{
		Name:  queueName,
		Model: &m,
	})

	worker.Pusher().SetDelay(delay).SetLevel(level).New(models)

	c.Ctx.WriteString("success")
}
