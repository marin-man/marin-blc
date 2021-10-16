package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
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