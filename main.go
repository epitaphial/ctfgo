package main

import (
	_"ctfgo/routers"
	"github.com/astaxie/beego"
	"ctfgo/tools"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.AddFuncMap("addone",tools.Addone)
	beego.Run()
}
