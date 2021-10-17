package main

import "bkc/BLC"

func main() {
	//bc := BLC.CreateBlockChain()
	//fmt.Println(bc.Tip)
	//bc.AddBlock([]byte("alice send 10 btc to bob"))
	//bc.AddBlock([]byte("bob send 5 btc to troytan"))
	//
	//bc.PrintChain()
	cli := BLC.CLI{}
	cli.Run()
}