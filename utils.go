package main

// 导入必要的模块
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// 定义一个int转为byte的函数，传入一个byte64的量，返回一串byte指针
func IntToByte(num int64) []byte {
	// 定义一个buffer空白的
	var buffer bytes.Buffer
	// 把数字写入buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	// 检查错误
	CheckError("检查错误() IntToByte() 行号: 18 utils.go", err)
	// 返回需要的buffer
	return buffer.Bytes()
}

// 定义一个检查错误的Error
func CheckError(pos string, err error){
	// 检查错误，如果出现错误，那么打印一些提示，并且退出程序
	if err != nil{
		fmt.Println("您的错误出现在: ", pos, "这是您的错误信息: ", err)
		// 退出程序 return code = 1
		os.Exit(1)
	}
}
// 定义两个辅助函数
func (input *TXInput) CanUnlockUTXOWith (unlockData string) bool{
	return input.ScriptSig == unlockData
}
func (output *TXOutput) CanUnlockedWith (unlockData string) bool{
	return output.ScriptPubKey == unlockData
}
// 定义一个Coinbase检查的函数
func (tx *Transaction) IsCoinBase() bool{
	if len(tx.TXInputs) == 1{
		if len(tx.TXInputs[0].TXID) == 0 && tx.TXInputs[0].Vout == -1{
			return true
		}
	}
	return false
}

// 定义一个拼接数据的函数
func (pow *ProofOfWork) PerpareData(nonce int64) []byte{
	// 计算一个Merkel Tree
	copy(pow.block.MerkelRoot, pow.block.TransactionHash())
	// 定义一个拼接数据的二位【】【】
	tmp := [][]byte{
		IntToByte(pow.block.Version),
		pow.block.PrevBlockHash,
		//pow.block.MerkelRoot = ,
		pow.block.MerkelRoot,
		IntToByte(pow.block.Timestamp),
		IntToByte(targetBits),
		IntToByte(nonce)}
	// 定义一个Join的函数
	data := bytes.Join(tmp, []byte{})
	// 返回数据
	return data
}
// 定义一个获取指定amout的UTXOs，内部我就不做注释了
func (bc *BlockChain) FindSuitableUTXOs(address string, amount float64) (map[string][]int64, float64){
	// 查找所有UTXO交易
	txs := bc.FindUTXOTransactions(address)
	// 定义一个交易大小
	var total float64
	// 所有可用的交易
	validUTXOs:= make(map[string][]int64)
	// 定义了一个名为FIND的跳转
FIND:
	for _, tx := range txs{
		// 定义TX输出
		outputs := tx.TXOutputs
		// 循环所有输出
		for index, output := range outputs{
			// 检查是否可以支配
			if output.CanUnlockedWith(address){
				// 可以支配
				// 检查余额
				if total < amount{
					total += output.Value
					validUTXOs[string(tx.TXID)] = append(validUTXOs[string(tx.TXID)], int64(index))
				}else {
					// 如果不可以支配，那么返回到FIND跳转点
					break FIND
				}
			}
		}
	}
	return validUTXOs, total
}
// 定义一个查找utxo交易的函数
func (bc *BlockChain)FindUTXOTransactions(address string) []Transaction{
	// 定义一些函数
	var UTXOTransactions []Transaction
	spentUTXO := make(map[string][]int64)
	// 定义一个迭代器，对区块链进行迭代
	it := bc.CreateIterator()
	// 循环迭代
	for {
		// 获取一个区块
		block := it.Next()
		// 进行循环判断
		for _, tx := range block.Transactions{
			if !tx.IsCoinBase() {
				for _, input := range tx.TXInputs {
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}
		OUTPUTS:
			for currIndex, output := range tx.TXOutputs{
				if spentUTXO[string(tx.TXID)] != nil{
					indexes := spentUTXO[string(tx.TXID)]
					for _, index := range indexes{
						if int64(currIndex) == index{
							continue OUTPUTS
						}
					}
				}

				if output.CanUnlockedWith(address){
					UTXOTransactions = append(UTXOTransactions, *tx)
				}
			}
		}
		// 检查一些数据
		if len(block.PrevBlockHash) == 0{
			break
		}
	}
	// 返回
	return UTXOTransactions
}
// 查找所有可以支配的UTXO
func (bc * BlockChain) FindUTXO(address string) []TXOutput {
	// 定义一个输出交易的定义
	var UTXOs []TXOutput
	// 定义一个tt
	txs := bc.FindUTXOTransactions(address)
	// 循环处理
	for _, tx := range txs{
		for _, utxo := range tx.TXOutputs{
			if utxo.CanUnlockedWith(address){
				UTXOs = append(UTXOs, utxo)
			}
		}
	}
	// 返回UTXOs
	return UTXOs
}