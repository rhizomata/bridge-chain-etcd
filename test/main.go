package main

import (
	"os"

	"github.com/rhizomata/bridge-chain-etcd/protocol"
)

func main() {
	args := os.Args[1:]

	client := protocol.NewClient("http://127.0.0.1:8080")

	if len(args) > 0 {
		if args[0] == "remove" {
			jobid := args[1]
			client.RemoveJob(jobid)
		} else if args[0] == "add" {
			client.AddJob([]byte(`#eth-relay:{"handler":"erc20","cas":[
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xB8c77482e45F1F44dE1745F52C74426C631bDD52"]}`))

			client.AddJob([]byte(`#eth-relay:{"handler":"erc721","cas":[
				"0x0e3a2a1f2146d86a604adc220b4967a898d7fe07",
				"0x06012c8cf97bead5deae237070f9587f8e7a266d"]}`))
		}
	}

}
