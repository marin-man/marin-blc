package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// IntToHex 实现 int64 转成 []byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if nil != err {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// JSONToSlice 标准 JSON 格式转切片
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	// json
	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}


// IsValidArgs 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		// 直接退出
		os.Exit(1)
	}
}

// Reverse 反转切片
func Reverse(data []byte)  {
	for i, j := 0, len(data) - 1; i < j; i, j = i + 1, j - 1 {
		data[i], data[j] = data[j], data[i]
	}
}

// GetEnvNodeId 获取节点 ID
func GetEnvNodeId() string {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	return nodeID
}