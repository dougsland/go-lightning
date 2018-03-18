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
	if argCount < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s bolt11\n", os.Args[0])
		os.Exit(1)
	}

	bolt11 := os.Args[1]

	data := lightning.JsonDecodePay{}
	ln := lightning.LightningRpc(common.RPC_FILENAME)

	json_data := ln.Decodepay(bolt11)

	err := json.Unmarshal([]byte(json_data), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Payment decoded")
	fmt.Println("======================")

	fmt.Println("Currency:", data.Result.Currency)
	fmt.Println("Timestamp:", data.Result.Timestamp)
	fmt.Println("Created_at:", data.Result.Created_at)
	fmt.Println("Expiry:", data.Result.Expiry)
	fmt.Println("Payee:", data.Result.Payee)
	fmt.Println("Msatoshi:", data.Result.Msatoshi)
	fmt.Println("Description:", data.Result.Description)
	fmt.Println("Min final cltv expiry:", data.Result.Min_final_cltv_expiry)
	fmt.Println("Payment hash:", data.Result.Payment_hash)
	fmt.Println("Signature:", data.Result.Signature)
	fmt.Println("--------------------------------------------------")
	ln.Destroy()

	os.Exit(0)
}
