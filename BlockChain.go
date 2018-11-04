package main

import (
	"STCBlockChain/bolt"
	"fmt"
	"os"
)

// 定义一些程序参数
const DbFile = "./blockchain.db" // 定义DbFile进行数据保存
const BucketName = "bcc" // 定义blockchain的Bucket bolt名称
const lastHashKey = "tailHash" // 尾巴hash的key，在数据库中存储
const genesisInfo = "讨厌的赵雅宁" // 这个是创世区块的内容，我写了一个我的坏班长

// 定义区块链类
type BlockChain struct {
	// 定义数据库操作句柄
	db   *bolt.DB
	// 尾巴hash的内容
	tail []byte
}

// 定义区块链迭代器，进行区块链的迭代，不断从数据库中获取Iter
type BlockChainIterator struct {
	// 现在的hash值
	currHash []byte
	// 数据库操作句柄
	db       *bolt.DB
}

// 检查数据库是否被创建，返回一个bool
func isDBExist() bool {
	// 检查Stat
	_, err := os.Stat(DbFile)
	// 检查数据库文件位置，是否有数据库文件
	if os.IsNotExist(err) {
		// 如果有，那么返回false
		return false
	}
	// 如果没有，那么返回true
	return true
}

// 定义创建区块链的函数
func CreateBlockChain(address string) *BlockChain {
	// 打开我们的数据库文件，mode操作符是0600
	db, err := bolt.Open(DbFile, 0600, nil)
	// 检查数据库返回
	CheckError("检查数据库() bolt.Open() BlockChain.go 第49行", err)
	// 定义尾巴hash
	var lastHash []byte
	// 打开一个数据库事务
	db.Update(func(tx *bolt.Tx) error {
		// 获取一个bucket，检查bucket
		bucket := tx.Bucket([]byte(BucketName))
		// 检查bucket是否为nil，如果为nil，那么我们直接进行创建，如果不为nil，获取lastHash，尾巴hash
		if bucket != nil {
			// 获取尾巴hash
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			// 创建一个挖矿交易，然后传入一些创世快的信息
			coinbase := CreateCoinbaseTx(address, genesisInfo)
			// 构建一个创世快
			genesis := CreateGenesisBlock(coinbase)
			// 创建一个bucket
			bucket, err := tx.CreateBucket([]byte(BucketName))
			// 检查错误
			CheckError("创建Bucket检查 BlockChain.go 67行", err)
			// 写入一些数据，把块进行序列化
			bucket.Put(genesis.Hash, genesis.Serialize())
			bucket.Put([]byte(lastHashKey), genesis.Hash)
			// 然后把hash进行赋值，把尾巴hash赋值为genesis.Hash
			lastHash = genesis.Hash
		}
		return nil
	})
	// 返回区块链的指针
	return &BlockChain{db: db, tail: lastHash}
}

// 获取区块链返回hadler
func GetBlockChainHandler() *BlockChain {
	// 检查是否拥有区块链，如果没有，直接退出
	if !isDBExist() {
		// 显示错误信息并退出
		fmt.Println("没有区块链")
		os.Exit(1)
	}
	// 打开一个数据库文件模式 0600
	db, err := bolt.Open(DbFile, 0600, nil)
	// 检查错误
	CheckError("检查数据库错误 bolt.Open() BlockChain.go 92行", err)
	// 尾巴hash
	var lastHash []byte
	// 打开一个查询事务
	db.View(func(tx *bolt.Tx) error {
		// 打开一个bucket，并且检查bucket
		bucket := tx.Bucket([]byte(BucketName))
		// 检查bucket是否为空，如果为空，那么需要直接退出，如果不为空，那么取出尾巴hash
		if bucket != nil {
			lastHash = bucket.Get([]byte(lastHashKey))
		}
		return nil
	})
	// 返回区块链的指针
	return &BlockChain{db: db, tail: lastHash}
}

// 添加一个块，传入一个交易
func (bc *BlockChain) AddBlock(txs[]*Transaction) {
	// 把上一个区块的hash取出
	var prevBlockHash []byte
	// 定义bc.db.View进行数据库的查询，打开查询事务
	bc.db.View(func(tx *bolt.Tx) error {
		// 打开一个bucket，进行操作
		bucket := tx.Bucket([]byte(BucketName))
		// 如果bucket为空，那么我们直接退出
		if bucket == nil {
			os.Exit(1)
		}
		// 获取上一级的hash
		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	// 创建一个block，传入一个挖矿交易和上一个hash
	block := CreateBlock(txs,prevBlockHash)
	// 打开一个数据库更新
	bc.db.Update(func(tx *bolt.Tx) error {
		// 取出一个bucket
		bucket := tx.Bucket([]byte(BucketName))
		// 写入一些数据
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(lastHashKey), block.Hash)
		// 获取尾巴
		bc.tail = block.Hash
		return nil
	})
}
// 定义一个创建迭代器的函数，这个函数是在BlockChainIterator中的函数
func (bc *BlockChain) CreateIterator() *BlockChainIterator {
	// 返回一个BlockChain迭代器的指针
	return &BlockChainIterator{currHash: bc.tail, db: bc.db}
}
// 定义一个Next()函数在BlockChainIter中，返回一个block
func (it *BlockChainIterator) Next() (block *Block) {
	// 打开一个数据库事务
	err := it.db.View(func(tx *bolt.Tx) error {
		// 拿到一个bucket
		bucket := tx.Bucket([]byte(BucketName))
		// 如果bucket出现问题 == nil的情况下，返回一个nil
		if bucket == nil {
			return nil
		}
		// 获取数据 data获取一个currHash
		data := bucket.Get(it.currHash)
		// 反序列化block，把bytes变为一个类型
		block = DeSerialize(data)
		// 把现在的hash值变为这个区块的上一个区块的hash值
		it.currHash = block.PrevBlockHash
		return nil
	})
	// 检查数据库错误
	CheckError("数据库检查() it.db.View() 行数 147", err)
	// 返回所有(block *Block)定义的参数
	return
}
