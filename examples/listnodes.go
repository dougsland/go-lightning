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
		fmt.Println("Nodeid:", element.Nodeid)
		fmt.Println("Alias:", element.Alias)
		fmt.Println("Color:", element.Color)
		fmt.Println("Last Timestamp:", element.Last_Timestamp)
		for _, net_element := range element.Addresses {
			fmt.Println("Type Address:", net_element.Type)
			fmt.Println("Address:", net_element.Address)
			fmt.Println("Port:", net_element.Port)
		}
		fmt.Println("--------------------------------------------------")
	}

	ln.Destroy()

	os.Exit(0)
}
