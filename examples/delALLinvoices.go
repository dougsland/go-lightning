/*
 * Example: List all invoices
 */
package main

import (
	"../lightning"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	var input string

	fmt.Print("Are you sure about DELETE ALL invoices? [Yes/No]")
	fmt.Scanln(&input)

	if input != "Yes" {
		fmt.Println("Aborting...")
		os.Exit(1)
	}

	data := lightning.JsonListInvoices{}
	ln := lightning.LightningRpc("/home/bitcoin/.lightning/lightning-rpc")

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

	if len(data.Result.Invoices) == 0 {
		fmt.Println("No invoices in the queue to remove")
		os.Exit(0)
	}
	for _, element := range data.Result.Invoices {
		fmt.Println("Removing label: ", element.Label)
		ln.Delinvoice(element.Label, element.Status)
	}

	ln.Destroy()

	os.Exit(0)
}
