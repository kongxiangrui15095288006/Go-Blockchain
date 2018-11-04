package main

// 导入必要的包
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// 定义一个名为Block的结构
type Block struct {
	// 定义版本
	Version int64
	// 定义上一个Block的值
	PrevBlockHash []byte
	// 定义Hash
	Hash []byte
	// 定义Merkel根
	MerkelRoot []byte
	// 定义时间戳
	Timestamp int64
	// 定义难度
	Bits int64
	// 定义随机值
	Nonce int64

	// 定义交易内容
	Transactions []*Transaction
}

// 定义一个创建区块的函数，传入一个交易 Transcation 和一个上区块哈希值返回一个区块（Block）类型的变量
func CreateBlock(txs []*Transaction, prevBlockHash []byte) *Block{
	// 定义一个区块类型的
	var block Block
	// 创建一个Block的结构体
	block = Block{
		Version: 1,
		PrevBlockHash:prevBlockHash,
		MerkelRoot: []byte{},
		Timestamp:time.Now().Unix(),
		Bits:targetBits,
		Nonce: 0,
		Transactions: txs}
	// 创建工作量证明进行挖矿
	pow := CreateProofOfWork(&block)
	// 进行挖矿，需要一些计算资源进行挖矿
	nonce, hash := pow.RunMine()
	// 挖矿完毕后，将挖出的Hash和Nonce赋值到block中，完成block的创建
	block.Nonce = nonce
	block.Hash = hash
	// 返回我们创建完成的block
	return &block
}

// 定义一个序列化数据的函数，把输入的block序列化为一堆bytes
func (block *Block) Serialize() []byte{
	// 创建一个空白的buffer用于存储以前的bytes
	var buffer bytes.Buffer
	// 新建一个encoder
	encoder := gob.NewEncoder(&buffer)
	// 对数据进行序列化
	err := encoder.Encode(block)
	// 检查出错
	CheckError("序列化 Serialize()出错提示 Block.go 65行", err)
	// 返回我们需要的buffer
	return buffer.Bytes()
}

// 定义一个反序列化的函数，把输入的byte数据序列化为一段Block类型，直接进行操作
func DeSerialize(data []byte) *Block{
	// 检查数据长度 data
	if len(data) == 0{
		// 返回并且值为nil
		return nil
	}
	// 定义一个block的数据
	var block Block
	// 定义decoder的数据并进行反序列化得到我们的decoder
	decoder := gob.NewDecoder(bytes.NewReader(data))
	// 定义decode函数并进行decode，把数据放入block容器指针中
	err := decoder.Decode(&block)
	// 检查错误
	CheckError("反序列化 DeSerialize()出错提示 Block.go 84行", err)
	// 返回block类型的数据
	return &block
}

// 创建一个创世快，也是区块链的头快，需要传入一个coinbase挖矿请求的交易，返回一个coinbase的快
func CreateGenesisBlock(coinbase *Transaction) *Block {
	// 返回coinbase块，并且创建一个数组，他是一个内部只有一个coinbase交易的块
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// 定义一个计算交易hash值得函数，这个函数可以计算交易hash值并进行返回一个 []byte类型的Go切片
func (block *Block)TransactionHash() []byte{
	// 定义一个拼接类型的byte[][]二位切片
	var TXHashes [][]byte
	// 获取一个交易
	txs := block.Transactions
	// 循环交易并进行拼接，大部分是一个交易
	for _, tx := range txs{
		// 使用txhashes进行拼接，如果有N个交易，那么拼接N次
		TXHashes = append(TXHashes, tx.TXID)
	}
	// 拼接数据，使用byteJoin进行拼接，把二维的数据拼接为1dim的切片
	data := bytes.Join(TXHashes, []byte{})
	// 计算hash，并进行Sum256
	hash := sha256.Sum256(data)
	// 返回hash[:]
	return hash[:]
}