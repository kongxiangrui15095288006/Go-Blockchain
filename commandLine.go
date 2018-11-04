package main

import (
	"flag"
	"fmt"
	"os"
)

const usage  = `
	printchain 输出区块链的内容 ： printchain
	createaccount 创建新的address区块链 : createaccount --address ADDRESS
	getbalance 获取一个addres下的余额 ： getbalance --address ADDRESS
	send 发送N个XTC到其他的账户中 ：send --from FROM ADDRESS --to TO ADDRESS --amount AMOUNT
	mining 挖矿，挖出一些XTC ：mining
`
const PrintChainCmdString = "printchain"
const CreateChainCmdString = "creataccount"
const getBalanceCmdString = "getbalance"
const sendCmdString = "send"
const miningCmdString = "mining"

type CLI struct {}

func (cli *CLI) PrintUsage() {
	fmt.Println("Invalid Input!!+")
	fmt.Println(usage)
	os.Exit(1)
}

func (cli *CLI) parameterCheck() {
	if len(os.Args) < 2{
		cli.PrintUsage()
	}
}

func (cli *CLI) Run() {
	cli.parameterCheck()

	// addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(getBalanceCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(sendCmdString, flag.ExitOnError)
	miningCmd := flag.NewFlagSet(miningCmdString, flag.ExitOnError)

	// addBlocCmdPara := addBlockCmd.String("data", "", "交易信息")
	CreateChainCmdPara := createChainCmd.String("address", "", "交易地址")
	GetBalanceCmdPara := getBalanceCmd.String("address", "", "交易地址")
	fromPara := sendCmd.String("from", "", "从那个账户")
	toPara := sendCmd.String("to", "", "去哪里")
	amountPara := sendCmd.Float64("amount", 0, "钱数")
	miningCmdPara := miningCmd.String("address", "", "钱包")

	switch os.Args[1] {
	case getBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckError("aaa()", err)
		if getBalanceCmd.Parsed() {
			if *GetBalanceCmdPara == ""{
				cli.PrintUsage()
			}
			cli.GetBalance(*GetBalanceCmdPara)
		}
	case sendCmdString:
		err := sendCmd.Parse(os.Args[2:])
		CheckError("aaa()", err)
		if sendCmd.Parsed(){
			if *fromPara == "" && *toPara == "" && *amountPara == 0{
				cli.PrintUsage()
			}
			cli.Send(*fromPara, *toPara, *amountPara)
		}
	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckError("sss()", err)
		if printChainCmd.Parsed(){
			cli.printChain()
		}
	case CreateChainCmdString:
		err := createChainCmd.Parse(os.Args[2:])
		CheckError("aaa()", err)
		if createChainCmd.Parsed(){
			cli.CreateChain(*CreateChainCmdPara)
		}
	case miningCmdString:
		err := miningCmd.Parse(os.Args[2:])
		CheckError("aaa()", err)
		if miningCmd.Parsed(){
			cli.Mining(*miningCmdPara)
		}
	default:
		cli.PrintUsage()
	}

}
