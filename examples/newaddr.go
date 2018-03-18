/*
 * Example: Get information from the system
 */
package main

import (
	"../lightning"
	"./common"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	data := lightning.JsonNewaddr{}
	typeaddr := "p2sh-segwit"
	if len(os.Args[1:]) == 1 {
		typeaddr = os.Args[1]
	}

	ln := lightning.LightningRpc(common.RPC_FILENAME)
	myaddr := ln.Newaddr(typeaddr)

	err := json.Unmarshal([]byte(myaddr), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(data.Result.Address) == 0 {
		fmt.Print("Unable to proceed, ",
			"JSON error message: ", myaddr)
		os.Exit(1)
	}

	fmt.Println("Information")
	fmt.Println("======================")
	fmt.Println("Type:", typeaddr)
	fmt.Println("Address:", data.Result.Address)

	ln.Destroy()

	os.Exit(0)
}
