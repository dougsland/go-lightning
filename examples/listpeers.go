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
	data := lightning.JsonListPeers{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds
	json_data := ln.Listpeers()

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Peers")
	fmt.Println("======================")

	peers := 0
	for _, element := range data.Result.Peers {
		fmt.Println("State:", element.State)
		fmt.Println("Id:", element.Id)
		fmt.Println("Netaddr:", element.Netaddr)
		fmt.Println("Connected:", element.Connected)
		fmt.Println("Alias:", element.Alias)
		fmt.Println("Color:", element.Color)
		fmt.Println("Owner:", element.Owner)
		fmt.Println("Id:", data.Id)
		peers++
		fmt.Println("--------------------------------------------------")
	}
	fmt.Println("Total peers:", peers)
	ln.Destroy()

	os.Exit(0)
}
