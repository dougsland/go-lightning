/*
 * Example: Fundchannel
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
	if argCount < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s channel_id satoshi\n", os.Args[0])
		os.Exit(1)
	}

	id := os.Args[1]
	satoshi := os.Args[2]

	ln := lightning.LightningRpc(common.RPC_FILENAME)

	ln.Fundchannel(id, satoshi)

	ln.Destroy()

	os.Exit(0)
}
