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
	data := lightning.JsonChannels{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds

	json_data := ""
	if len(os.Args) > 1 {
		json_data = ln.Listchannels(os.Args[1])
	} else {
		json_data = ln.Listchannels()

	}

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Channels")
	fmt.Println("======================")

	for _, element := range data.Result.Channels {
		fmt.Println("Source:", element.Source)
		fmt.Println("Destination:", element.Destination)
		fmt.Println("Short_channel_id:", element.Short_channel_id)
		fmt.Println("Flags:", element.Flags)
		fmt.Println("Active:", element.Active)
		fmt.Println("Public:", element.Public)
		fmt.Println("Last_update:", element.Last_update)
		fmt.Println("Base_fee_millisatoshi:", element.Base_fee_millisatoshi)
		fmt.Println("Fee_per_millionth:", element.Fee_per_millionth)
		fmt.Println("Delay:", element.Delay)
		fmt.Println("--------------------------------------------------")
	}
	ln.Destroy()

	os.Exit(0)
}
