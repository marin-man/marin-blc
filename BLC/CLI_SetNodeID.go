package BLC

import (
	"log"
	"os"
)

// SetNodeId 设置端口号（环境变量）
func (cli *CLI) SetNodeId(nodeId string) {
	err := os.Setenv("NODE_ID", nodeId)
	if nil != err {
		log.Fatalf("set env failed! %v\n", err)
	}
}
