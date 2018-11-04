package main

import (
	"fmt"
)

func (cli *CLI) printChain(){
	bc := GetBlockChainHandler()
	it := bc.CreateIterator()

	for {
		block := it.Next()
		fmt.Printf("版本: %d\n", block.Version)
		fmt.Printf("上区块哈希值：%x\n", block.PrevBlockHash)
		fmt.Printf("哈希值：%x\n", block.Hash)
		fmt.Printf("时间戳：%d\n", block.Timestamp)
		fmt.Printf("难度：%d\n", block.Bits)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("数据：%v\n", block.Transactions)

		if len(block.PrevBlockHash) == 0{
			break
		}
	}
}


func (cli *CLI) CreateChain(address string) {
	bc := CreateBlockChain(address)
	defer bc.db.Close()
	fmt.Printf("区块链创建成功!")
}

func (cli *CLI) GetBalance(address string){
	bc := GetBlockChainHandler()
	utxos := bc.FindUTXO(address)
	var total float64 = 0
	for _, utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("余额: %f", total)
}

func (cli *CLI) Send(from string, to string, amount float64){
	bc := GetBlockChainHandler()
	tx := NewTransction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
}

func (cli *CLI) Mining(address string){
	bc := GetBlockChainHandler()
	for {
		coinbase := CreateCoinbaseTx(address, genesisInfo)
		bc.AddBlock([]*Transaction{coinbase})
	}
}