/*
 * Example: List all invoices
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
	data := lightning.JsonListInvoices{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds
	json_data := ln.Listinvoices()

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Invoices")
	fmt.Println("======================")
	for _, element := range data.Result.Invoices {
		fmt.Println("Label:", element.Label)
		fmt.Println("Payment Hash:", element.Payment_hash)
		fmt.Println("Msatoshi:", element.Msatoshi)
		fmt.Println("Status:", element.Status)
		fmt.Println("Expiry time:", element.Expiry_time)
		fmt.Println("Expires at:", element.Expires_at)
		fmt.Println("--------------------------------------------------")
	}

	ln.Destroy()

	os.Exit(0)
}
