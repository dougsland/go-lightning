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

	argCount := len(os.Args[1:])
	if argCount < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s label status\n", os.Args[0])
		os.Exit(1)
	}

	lb := os.Args[1]
	status := os.Args[2]

	data := lightning.JsonInvoice{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	myinvoice := ln.Delinvoice(lb, status)

	err := json.Unmarshal([]byte(myinvoice), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(data.Result.Label) == 0 {
		fmt.Print("Unable to proceed with invoice, ",
			"JSON error message: ", myinvoice)
		os.Exit(1)
	}

	fmt.Println("Removed Invoice")
	fmt.Println("======================")
	fmt.Println("Label:", data.Result.Label)
	fmt.Println("Payment hash:", data.Result.Payment_hash)
	fmt.Println("Msatoshi:", data.Result.Msatoshi)
	fmt.Println("Expiry time:", data.Result.Expiry_time)
	fmt.Println("Expires at:", data.Result.Expires_at)
	fmt.Println("Id:", data.Id)
	ln.Destroy()

	os.Exit(0)
}
