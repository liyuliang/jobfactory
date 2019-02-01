package controllers

import (
	"github.com/astaxie/beego"
	"github.com/liyuliang/models/protobuf"
	"jobfactory/worker"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	url := c.GetString("url")
	queueName := c.GetString("queue")

	if url == "" {
		c.Abort("404")
		return
	}

	models := []*worker.Model{}

	m := protobuf.ParserManhuaPage{
		Url: url,
	}

	models = append(models, &worker.Model{
		Name:  queueName,
		Model: &m,
	})

	worker.Pusher().New(models)

	c.Ctx.WriteString("success")
}
