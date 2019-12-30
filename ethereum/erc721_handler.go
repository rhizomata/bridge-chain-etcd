package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rhizomata/bridge-chain-etcd/kernel/worker"
)

// ERC721LogHandler implements LogHandler
type ERC721LogHandler struct{}

// Name : erc20
func (handler *ERC721LogHandler) Name() string { return "erc721" }

// HandleLog ..
func (handler *ERC721LogHandler) HandleLog(helper *worker.Helper, log types.Log) error {
	address := log.Address.Hex()
	logHash := log.Topics[0].Hex()

	if len(log.Topics) > 2 {
		fromAddr := log.Topics[1].Hex()
		toAddr := log.Topics[2].Hex()
		fmt.Println("LOG ", handler.Name(), "- addr:", address, ", loghash:", logHash, ",fromAddr=", fromAddr, ",toAddr=", toAddr)
	} else {
		fmt.Println("LOG ", handler.Name(), "- addr:", address, ", loghash:", logHash)
	}

	return nil
}
