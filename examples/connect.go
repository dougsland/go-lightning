/*
 * Example: Connect to node
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
	if argCount < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s nodeid host port\n", os.Args[0])
		os.Exit(1)
	}

	nodeid := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	ln := lightning.LightningRpc(common.RPC_FILENAME)

	ln.Connect(nodeid, host, port)

	ln.Destroy()
	os.Exit(0)
}
