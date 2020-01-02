package ethereum

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rhizomata/bridge-chain-etcd/kernel/worker"
)

// EthSubscriber implements worker.Worker
type EthSubscriber struct {
	id         string
	client     *ethclient.Client
	networkURL string
	jobInfo    *EthSubsJobInfo
	helper     *worker.Helper
	started    bool
	handler    LogHandler
}

// LogHandler ..
type LogHandler interface {
	Name() string
	HandleLog(helper *worker.Helper, log types.Log) error
}

// EthSubsJobInfo ..
type EthSubsJobInfo struct {
	Handler           string   `json:"handler"`
	CAs               []string `json:"cas"`
	contractAddresses []common.Address
	From              uint64 `json:"from"`
}

// EthSubsManager implements worker.Factory, name eth_subs
type EthSubsManager struct {
	networkURL string
	handlers   map[string]LogHandler
}

// BlockCheckPoint ..
type BlockCheckPoint struct {
	BlockNumber uint64
	Index       uint
}

// var erc20abi = string(`
// [{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"tokenOwner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"}]
// `)

// NewEthSubsManager ..
func NewEthSubsManager(networkURL string) *EthSubsManager {
	manager := EthSubsManager{networkURL: networkURL}
	// contractAbi, err := abi.JSON(strings.NewReader(erc20abi))
	// if err != nil {
	// 	log.Fatal("Load ABI : ", err)
	// }
	// manager.contractAbi = contractAbi
	manager.handlers = make(map[string]LogHandler)

	// Register default built-in handlers
	manager.RegisterLogHandler(&ERC20LogHandler{})
	manager.RegisterLogHandler(&ERC721LogHandler{})

	return &manager
}

// RegisterLogHandler ..
func (manager *EthSubsManager) RegisterLogHandler(handler LogHandler) {
	manager.handlers[handler.Name()] = handler
}

// Name implements worker.Factory.Name
func (manager *EthSubsManager) Name() string {
	return "eth_subs"
}

// NewWorker implements worker.Factory.NewWorker
func (manager *EthSubsManager) NewWorker(helper *worker.Helper) (wroker worker.Worker, err error) {
	jobInfo := new(EthSubsJobInfo)
	json.Unmarshal(helper.Job(), jobInfo)

	addrs := []common.Address{}

	for _, ca := range jobInfo.CAs {
		addr := common.HexToAddress(ca)
		addrs = append(addrs, addr)
	}

	jobInfo.contractAddresses = addrs

	handler := manager.handlers[jobInfo.Handler]

	if handler == nil {
		log.Println("[ERROR-WorkerMan] Unknown Log Handler ", jobInfo.Handler)
		return nil, errors.New("Unknown Log Handler " + jobInfo.Handler)
	}

	subscriber := EthSubscriber{id: helper.ID(), jobInfo: jobInfo, networkURL: manager.networkURL,
		helper: helper, handler: handler}

	return &subscriber, nil
}

//ID ..
func (subscriber *EthSubscriber) ID() string {
	return subscriber.id
}

//Start ..
func (subscriber *EthSubscriber) Start() error {
	if subscriber.client != nil {
		subscriber.client.Close()
	}
	client, err := ethclient.Dial(subscriber.networkURL)
	if err != nil {
		log.Fatal("[ERROR] Cannot Connect to ", subscriber.networkURL, err)
	}
	subscriber.client = client

	log.Println("[Debug] ETH Subs :", subscriber.jobInfo.CAs, ", from:", subscriber.jobInfo.From)
	checkPoint := &BlockCheckPoint{}
	subscriber.helper.GetCheckpoint(checkPoint)

	go func() {
		if checkPoint.BlockNumber > 0 {
			subscriber.collect(checkPoint)
		}

		subscriber.subscribe(checkPoint)
		log.Println("[WARN] ETH Subs Ends. ", subscriber.ID())
	}()
	return nil
}

func (subscriber *EthSubscriber) handleLog(elog types.Log, checkPoint *BlockCheckPoint) {
	err := subscriber.handler.HandleLog(subscriber.helper, elog)
	if err != nil {
		log.Println("[FATAL-ETH-LogHandler] ", subscriber.ID(), err)
	}
	checkPoint.BlockNumber = elog.BlockNumber
	checkPoint.Index = elog.Index
	subscriber.helper.PutCheckpoint(checkPoint)
}

func (subscriber *EthSubscriber) subscribe(checkPoint *BlockCheckPoint) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber)),
		Addresses: subscriber.jobInfo.contractAddresses,
	}

	logs := make(chan types.Log)
	sub, err := subscriber.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Println("[ERROR] SubscribeFilterLogs ", subscriber.ID(), err)
		return
	}

	defer sub.Unsubscribe()

	subscriber.started = true

	for subscriber.started {
		select {
		case err := <-sub.Err():
			if !subscriber.started {
				break
			}
			log.Print("[ERROR] Eth Sub ", subscriber.ID(), err)
		case vLog := <-logs:
			if !subscriber.started {
				log.Println("[WARN] Eth Subscriber Stops .. ", subscriber.ID())
				break
			}

			// fmt.Printf("Sub Log Block Number: %d:%d  Addr: %s\n", vLog.BlockNumber, vLog.Index, vLog.Address.Hex())

			subscriber.handleLog(vLog, checkPoint)
		}
	}

	subscriber.started = false
}
func (subscriber *EthSubscriber) collect(checkPoint *BlockCheckPoint) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber)),
		Addresses: subscriber.jobInfo.contractAddresses,
	}

	logs, err := subscriber.client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Println("[ERROR-ETH-Subs]", err)
		return
	}

	for _, vLog := range logs {
		if vLog.BlockNumber == checkPoint.BlockNumber && vLog.Index <= checkPoint.Index {
			// fmt.Println("------ Skip Handle Log : Block - ", vLog.BlockNumber, ", Index - ", vLog.Index, "<=", checkPoint.Index)
			continue
		}

		log.Printf("Collect Log - %d:%d \n", vLog.BlockNumber, vLog.Index)
		subscriber.handleLog(vLog, checkPoint)
	}
}

//Stop ..
func (subscriber *EthSubscriber) Stop() error {
	if subscriber.client != nil {
		subscriber.client.Close()
	}

	subscriber.started = false
	return nil
}

//IsStarted ..
func (subscriber *EthSubscriber) IsStarted() bool {
	return subscriber.started
}
