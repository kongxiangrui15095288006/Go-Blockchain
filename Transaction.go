package main

// 导入必要的包
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// 定义一个挖矿reward也就是挖矿奖励
const reward = 12.5;

// 定义一个交易结构
type Transaction struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput
}

// 定义一个交易输入结构
type TXInput struct {
	TXID []byte
	Vout int64
	ScriptSig string
}

// 定义一个交易输出结构
type TXOutput struct {
	Value float64
	ScriptPubKey string
}

// 设置交易ID
func (tx *Transaction) SetTXID() {
	// 定义一个空白的buffer
	var buffer bytes.Buffer
	// 定义一个encoder
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	// 检查错误
	CheckError("序列化数据() NewEncoder() Transaction.go 40行", err)
	// 计算hash
	hash := sha256.Sum256(buffer.Bytes())
	// 把txID设为hash的切片
	tx.TXID = hash[:]
}

// 创建一个Coinbase交易
func CreateCoinbaseTx(address string, data string) *Transaction{
	// 如果没有传入挖矿交易的数据，那么需要进行处理
	if data == ""{
		data = "2018年11月3日一个小学生所写！"
	}
	// 构建一个挖矿交易
	inputs := TXInput{[]byte{}, -1, data}
	outputs := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{inputs}, []TXOutput{outputs}}
	tx.SetTXID()
	// 返回他
	return &tx
}

// 定义一个新的交易
func NewTransction(from, to string, amount float64, bc *BlockChain) *Transaction{
	// 定义一些交易数据
	var vaildUTXOs = make(map[string][]int64)
	var total float64
	var inputs []TXInput
	var outputs []TXOutput
	// 查找UTXOs
	vaildUTXOs, total = bc.FindSuitableUTXOs(from, amount)
	// 检查余额
	if amount > total{
		// 输出信息
		fmt.Println("您无法进行汇款，因为您的余额已不足！")
		// 直接退出程序
		return nil
	}
	// 进行循环处理
	for txId, outputIndexes := range vaildUTXOs{
		for _,index :=range outputIndexes{
			input := TXInput{[]byte(txId), int64(index), from}
			inputs = append(inputs, input)
		}
	}
	// 构建一些数据
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	// 找零钱
	if total > amount {
		output := TXOutput{total-amount, from}
		outputs = append(outputs, output)
	}
	// 构建交易类
	tx := Transaction{[]byte{}, inputs, outputs}
	// 返回数据
	return &tx
}