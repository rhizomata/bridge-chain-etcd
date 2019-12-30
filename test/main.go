package main

import "github.com/rhizomata/bridge-chain-etcd/protocol"

func main() {
	client := protocol.NewClient("http://127.0.0.1:8080")
	// client.AddJob([]byte(`{"ca":"0xdac17f958d2ee523a2206206994597c13d831ec7", "from":6383840}`))
	// client.AddJob([]byte(`{"ca":"0xB8c77482e45F1F44dE1745F52C74426C631bDD52"}`))
	// client.AddJob([]byte(`{"ca":"0xab95e915c123fded5bdfb6325e35ef5515f1ea69"}`))
	// client.AddJob([]byte(`{"ca":"0x6c37bf4f042712c978a73e3fd56d1f5738dd7c43"}`))
	// client.AddJob([]byte(`{"ca":"0x0e3a2a1f2146d86a604adc220b4967a898d7fe07"}`)) // 721 토큰

	client.AddJob([]byte(`{"handler":"erc20","cas":[
		"0xdac17f958d2ee523a2206206994597c13d831ec7","0xB8c77482e45F1F44dE1745F52C74426C631bDD52"]}`))

	client.AddJob([]byte(`{"handler":"erc721","cas":[
		"0x0e3a2a1f2146d86a604adc220b4967a898d7fe07","0x06012c8cf97bead5deae237070f9587f8e7a266d"]}`))
}
