package main

import (
	"log"

	"github.com/rhizomata/bridge-chain-etcd/api"
	"github.com/rhizomata/bridge-chain-etcd/ethereum"
	"github.com/rhizomata/bridge-chain-etcd/kernel"
	"github.com/rhizomata/bridge-chain-etcd/kernel/cluster"
	"github.com/rhizomata/bridge-chain-etcd/kernel/job"
	"github.com/rhizomata/bridge-chain-etcd/kernel/model"
	"github.com/rhizomata/bridge-chain-etcd/protocol"
)

func main() {
	daemonConfig := model.ParseFlagConfig()
	daemonAddr := daemonConfig.GetDaemonAddr()

	// "wss://mainnet.infura.io/ws"
	tokenSubsMan := ethereum.NewEthSubsManager("wss://mainnet.infura.io/ws")
	kernel := kernel.New(daemonConfig, tokenSubsMan.NewWorker)

	kernel.SetJobOrganizer(job.NewSimpleOrganizer())

	kernel.GetClusterManager().SetHealthCheckDelegator(func(memb *cluster.Member) bool {
		return protocol.CheckHealth(memb.DaemonURL)
	})

	err := kernel.Start()
	if err != nil {
		log.Fatal("[ERROR] Daemon Start Fail", err)
	}

	apiServer := api.StartServer(kernel, daemonAddr)

	<-apiServer.Error()
}
