package BLC

// GetData 请求指定区块
type GetData struct {
	AddrFrom	string		// 当前地址
	ID			[]byte		// 区块哈希
}