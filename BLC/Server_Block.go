package BLC

type BlockData struct {
	AddrFrom 	string		// 节点地址
	Block		[]byte		// 区块数据（序列化数据）
}