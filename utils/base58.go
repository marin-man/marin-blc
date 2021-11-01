package utils

import (
	"bytes"
	"math/big"
)

// base58 编码

// base58 编码基数表
var b58Alphabet = []byte("" +
	"123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")

// Base58Encode 编码函数
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	// 求余的基本长度
	base := big.NewInt(int64(len(b58Alphabet)))
	// 求余数和商
	zero := big.NewInt(0)
	// 设置余数，代表 base58 基数表的索引位置
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}
	// 反转 result切片
	Reverse(result)
	// 添加一个前缀 1，代表是一个地址
	result = append([]byte{b58Alphabet[0]}, result...)
	return result
}

// Base58Decode 解码函数
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 1
	// 去掉前缀
	data := input[zeroBytes:]
	for _, b := range data {
		// 查找 input 中指定数字/字符在基数表中出现的索引
		charIndex := bytes.IndexByte(b58Alphabet, b)
		// 余数 * 58
		result.Mul(result, big.NewInt(int64(len(b58Alphabet))))
		// 乘积结果 + mod
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	// 转换为 byte 字节数组
	decoded := result.Bytes()
	return decoded
}