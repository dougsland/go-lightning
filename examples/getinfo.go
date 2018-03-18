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
	data := lightning.JsonGetInfo{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	json_data := ln.Getinfo()

	json.Unmarshal([]byte(json_data), &data)

	fmt.Println("Information")
	fmt.Println("======================")
	fmt.Println("Id:", data.Result.Id)
	fmt.Println("Port:", data.Result.Port)

	if len(data.Result.Address) > 0 {
		for _, net_element := range data.Result.Address {
			fmt.Println("Type Address:", net_element.Type)
			fmt.Println("Address:", net_element.Address)
			fmt.Println("Port:", net_element.Port)
		}
	}

	fmt.Println("Version:", data.Result.Version)
	fmt.Println("Blockheight:", data.Result.Blockheight)
	fmt.Println("Network:", data.Result.Network)

	ln.Destroy()

	os.Exit(0)
}
