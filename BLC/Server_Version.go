package BLC

// Version 当前区块版本信息（决定区块是否需要同步）
type Version struct {
	// Version		int		// 版本号
	Height		int		// 当前节点的区块高度
	AddrFrom	string	// 当前节点的地址
}