package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type CreateAccountController struct {
	beego.Controller
}
type GetBalanceController struct {
	beego.Controller
}
type SendController struct {
	beego.Controller
}
type MiningController struct {
	beego.Controller
}

func (this *CreateAccountController) Get() {
	Address := this.Input().Get("address")
	bc := CreateBlockChain(Address)
	fmt.Printf("区块链创建成功!\n")
	this.Ctx.WriteString("{'code': 0, 'status': true, 'message': '区块链创建成功'}")
	bc.db.Close()
}

func (this *GetBalanceController) Get() {
	Address := this.Input().Get("address")
	fmt.Println("ffff")
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	utxos := bc.FindUTXO(Address)
	var total float64 = 0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("余额: %f\n", total)
	this.Ctx.WriteString("{'code': 0, 'status': true, 'message': '获取余额成功', 'balance': " + strconv.FormatFloat(total, 'f', 1, 64) + "}")
}

func (this *SendController) Get() {
	From := this.Input().Get("from")
	To := this.Input().Get("to")
	Amount, _ := strconv.ParseFloat(this.Input().Get("amount"), 64)
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	tx := NewTransction(From, To, Amount, bc)
	bc.AddBlock([]*Transaction{tx})
	this.Ctx.WriteString("{'code': 0, 'status': true, 'message': '汇款成功，已经写入区块链'}")
}

func (this *MiningController) Get() {
	Address := this.Input().Get("address")
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	coinbase := CreateCoinbaseTx(Address, genesisInfo)
	bc.AddBlock([]*Transaction{coinbase})
	this.Ctx.WriteString("{'code': 0, 'status': true, 'message': '挖出了一个区块'}")
}