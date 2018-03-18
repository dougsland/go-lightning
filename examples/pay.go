/*
 * Example: Pay
 */
package main

import (
	"../lightning"
	"./common"
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

	ln := lightning.LightningRpc(common.RPC_FILENAME)

	// Default timeout for read operation in the API is 200 Milliseconds
	// Users can change it, example 300 Milliseconds
	ln.Readtimeout = 300 // Milliseconds

	print(ln.Pay(bolt11))

	ln.Destroy()

	os.Exit(0)
}
