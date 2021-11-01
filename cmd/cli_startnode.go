package cmd

import "bkc/network"

// startNode 节点启动服务
func (cli *CLI) startNode(nodeId string) {
	network.StartServer(nodeId)
}