package main

import (
	_ "jobfactory/routers"
	"github.com/astaxie/beego"
	"jobfactory/conf"
)

func main() {
	auth := beego.AppConfig.String("auth")
	conf.Remote().Load(auth)
	beego.Run()
}

