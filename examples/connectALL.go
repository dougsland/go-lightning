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
	"strconv"
)

func main() {
	data := lightning.JsonListNodes{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds
	json_data := ln.Listnodes()

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Nodes")
	fmt.Println("======================")

	for _, element := range data.Result.Nodes {
		for _, net_element := range element.Addresses {
			if len(net_element.Address) > 0 && net_element.Port == 9735 {
				net_port := strconv.FormatInt(int64(net_element.Port), 10)
				fmt.Println(
					fmt.Sprintf(
						`Trying to connect %s %s:%s`,
						element.Nodeid,
						net_element.Address,
						net_port,
					),
				)
				ln.Connect(element.Nodeid, net_element.Address, net_port)
			}
		}
	}

	ln.Destroy()

	os.Exit(0)
}
