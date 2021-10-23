package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

// 存入所有输出的集合

type TXOutputs struct {
	TXOutputs []*TxOutput
}

// Serialize 输出集合序列化
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(txOutputs); nil != err {
		log.Panicf("serialize the utxo failed! %v\n", err)
	}
	return result.Bytes()
}

// Deserializer 输出集合序列化
func Deserializer(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err := decoder.Decode(&txOutputsBytes); nil != err {
		log.Panicf("deserialize the struct utxo failed! %v\n")
	}
	return &txOutputs
}