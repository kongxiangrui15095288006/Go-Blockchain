package main

// 导入一些必要的包
import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// 设置难度
const targetBits = 4

// 定义一个工作量证明的结构
type ProofOfWork struct {
	// 定义内部有一个block
	block *Block

	// 定义内部有一个target为big.Int的函数
	target *big.Int
}

// 创建一个POW工作量证明的函数：
func CreateProofOfWork(block *Block) *ProofOfWork{
	// 定义一个target，新建一个BigInt类型的变量
	target := big.NewInt(1)
	// 定义一个Lsh
	target.Lsh(target, uint(256 - targetBits))
	// 新建一个POW结构
	pow := ProofOfWork{block: block, target: target}
	// 返回一个工作正明的结构
	return &pow
}

func (pow *ProofOfWork) RunMine() (int64, []byte) {
	// 定义一个随机数
	var nonce int64 = 0
	// 定义一个32位hash
	var hash [32]byte
	// 定义一个大int
	var hashInt big.Int
	// 输出一些信息
	fmt.Printf("正在挖矿 目标哈希值：%x 目标随机数： %d \n", pow.target.Bytes(), pow.block.Nonce)
	// 循环进行挖矿
	for nonce < math.MaxInt64{
		// 数据进行处理
		data := pow.PerpareData(nonce)
		// 计算hash值
		hash = sha256.Sum256(data)
		// 对hashint对setbytes
		hashInt.SetBytes(hash[:])
		// 进行随机数比对
		if hashInt.Cmp(pow.target) == -1{
			// 输出一些提示
			fmt.Printf("找到了一个区块 已经存入区块链\n")
			break
		}else{
			// 随机数加1
			nonce++
		}
	}
	// 返回一些数据
	return nonce, hash[:]
}