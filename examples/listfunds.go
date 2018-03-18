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
	data := lightning.JsonListFunds{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds
	json_data := ln.Listfunds()

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Outputs")
	fmt.Println("======================")

	for _, element := range data.Result.Outputs {
		fmt.Println("Txid:", element.Txid)
		fmt.Println("Output:", element.Output)
		fmt.Println("Value:", element.Value)
		fmt.Println("Status:", element.Status)
	}
	fmt.Println("--------------------------------------------------")

	for _, ch_element := range data.Result.Channels {
		fmt.Println("Peerid:", ch_element.Peer_id)
		fmt.Println("Short Channel id:", ch_element.Short_channel_id)
		fmt.Println("Channel Sat:", ch_element.Channel_sat)
		fmt.Println("Channel Total Sat:", ch_element.Channel_total_sat)
		fmt.Println("Funding Txid:", ch_element.Funding_txid)
	}
	fmt.Println("--------------------------------------------------")

	ln.Destroy()

	os.Exit(0)
}
