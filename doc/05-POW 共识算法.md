此前我们实现的区块可以计算 hash，但没有工作难度，只需一次计算就可以得到 hash 值，为了体现工作量，需要给矿工们增加挖矿难度。
```go
type Block struct {
    Timestamp       int64       // 时间戳
    Data            []byte      // 数据域
    PrevBlockHash   []byte      // 前块 hash 值
    Hash            []byte      // 当前块 hash 值
    Nonce           int64       // 随机值
}
```

## POW 算法
```go
var (
    // Nonce 循环上限
    maxNonce = math.MaxInt64
)

// 难度值
const targetBits = 24

// Pow 结构
type ProofOfWork struct {
    block   *Block
    target  *big.Int
}

// 创建 POW
func NewProofOfWork(b *Block) *ProofOfWork {
    // target 为最终难度值
    target := big.NewInt(1)
    // targt 为 1 向左位移 256-24（挖矿难度）
    target.Lsh(target, uint(256-targetBits))
    // 生成 pow 结构
    pow := &ProofOfWork{b, target}
    return pow
}

// 挖矿运行
func (pow *ProofOfWork) Run() (int, []byte) {
    var hashInt big.Int
    var hash [32]byte
    nonce := 0
    fmt.Printf("Mining the block containing %s, maxNonce = %d\n", pow.block.Data, maxNonce)
    for nonce < maxNonce {
        // 数据准备
        data := pow.prepareData(nonce)
        // 计算 hash
        hash = sha256.Sum256(data)
        hashInt.SetBytes(hash[:])
        // 按字节比较，小于 0 代表找到目标 Nonce
        if hashInt.Cmp(pow.target) == -1 {
            break
        } else {
            nonce++
        }
    }
    return nonce, hash[:]
}

// 准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
    data := bytes.Join(
        [][]byte{
            pow.block.PrevBlockHash,
            pow.block.Data,
            Int2Hex{pow.block.Timestamp},
            Int2Hex{int64{targetBits)),
            Int2Hex(nonce),
        },
        []byte{},
    )
    return data
}

// 将 int64 写入 []byte
func Int2Hex(num int64) []byte {
    buff := new(bytes.Buffer)
    // 大端法写入
    binary.Write(buff, binary.BigEndian, num)
    return buff.Bytes()
}

// 校验区块正确性
func (pow *ProofOfWork) Validate() bool {
    data := pow.prepareData(pow.block.Nonce)
    hashInt.SetBytes(hash[:])
    return hashInt.Cmp(pow.target) == -1
}
```

## 修改 block 代码，将 sethash 函数换掉
```go
// 创建区块
func NewBlock(data string, prevBlockHash, []byte) *Block {
    block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{})
    // 需要先挖矿
    pow := NewProofOfWork(block)
    nonce, hash := pow.Run()
    // 设置 hash 和 nonce
    block.Hash = hash
    block.Nonce = nonce
    return block
}
```