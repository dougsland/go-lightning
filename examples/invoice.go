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
	data := lightning.JsonInvoice{}

	argCount := len(os.Args[1:])
	if argCount < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s msatoshi label description\n", os.Args[0])
		os.Exit(1)
	}

	msatoshi, errUint := strconv.ParseUint(os.Args[1], 10, 64)
	println(msatoshi)
	if errUint != nil {
		fmt.Println(errUint)
		os.Exit(1)
	}

	label := os.Args[2]
	description := os.Args[3]

	ln := lightning.LightningRpc(common.RPC_FILENAME)

	myinvoice := ln.Invoice(msatoshi, label, description)

	err := json.Unmarshal([]byte(myinvoice), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(data.Result.Bolt11) == 0 {
		fmt.Print("Unable to proceed with invoice, ",
			"JSON error message: ", myinvoice)
		os.Exit(1)
	}

	fmt.Println("Invoice Information")
	fmt.Println("======================")
	fmt.Println("Bolt11:", data.Result.Bolt11)
	fmt.Println("Payment hash:", data.Result.Payment_hash)
	fmt.Println("Expiry time:", data.Result.Expiry_time)
	fmt.Println("Expires at:", data.Result.Expires_at)
	fmt.Println("Id:", data.Id)

	ln.Destroy()
	os.Exit(0)
}
