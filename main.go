package main

import "github.com/astaxie/beego"

func main() {
	beego.Router("/createaccount", &CreateAccountController{})
	beego.Router("/getbalance", &GetBalanceController{})
	beego.Router("/send", &SendController{})
	beego.Router("/mining", &MiningController{})
	beego.Run()
}
