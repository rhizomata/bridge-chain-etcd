package worker

import (
	"errors"
	"log"

	"github.com/rhizomata/bridge-chain-etcd/kernel/kv"
)

// Manager ..
type Manager struct {
	cluster             string
	localid             string
	kv                  kv.KV
	workerFactoryMethod func(helper *Helper) (Worker, error)
	workers             map[string]Worker
}

// NewManager create Manager
func NewManager(cluster string, localid string, kv kv.KV,
	workerFactoryMethod func(helper *Helper) (Worker, error)) *Manager {
	manager := Manager{cluster: cluster, localid: localid, kv: kv,
		workerFactoryMethod: workerFactoryMethod}
	manager.workers = make(map[string]Worker)
	return &manager
}

// Cluster get cluster name
func (manager *Manager) Cluster() string { return manager.cluster }

// LocalID get local kernel id
func (manager *Manager) LocalID() string { return manager.localid }

// KV get etcd kv
func (manager *Manager) KV() kv.KV { return manager.kv }

// RegisterWorker ..
func (manager *Manager) registerWorker(id string, job []byte) error {
	if manager.workers[id] != nil {
		return errors.New("Worker[" + id + "] is already registered. If you want register new one, DeregisterWorker first")
	}
	helper := NewHelper(manager.cluster, id, job, manager.kv)
	worker, err := manager.workerFactoryMethod(helper)
	if err != nil {
		log.Println("[ERROR] Cannot create worker ", err)
		return err
	}

	manager.workers[id] = worker
	err = worker.Start()
	return err
}

// DeregisterWorker ..
func (manager *Manager) deregisterWorker(id string) error {
	worker := manager.workers[id]
	if worker == nil {
		return errors.New("Worker[" + id + "] is not registered.")
	}

	err := worker.Stop()

	if err == nil {
		delete(manager.workers, id)
	}

	return err
}

// Dispose ..
func (manager *Manager) Dispose() error {
	array := []string{}
	for id := range manager.workers {
		array = append(array, id)
	}

	for _, id := range array {
		manager.deregisterWorker(id)
	}

	return nil
}

// SetJobs ..
func (manager *Manager) SetJobs(jobs map[string][]byte) {
	log.Println("[WARN-WorkerManager] Set Jobs:", len(jobs))

	tempWorkers := make(map[string]Worker)
	newWorkers := make(map[string]Worker)

	for id, worker := range manager.workers {
		tempWorkers[id] = worker
	}

	for id, data := range jobs {
		worker := tempWorkers[id]
		if worker != nil {
			delete(tempWorkers, id)
		} else {
			helper := NewHelper(manager.cluster, id, data, manager.kv) 
			// worker = manager.workerFactoryMethod(helper)
			worker2, err := manager.workerFactoryMethod(helper)
			if err != nil {
				log.Println("[ERROR-WorkerMan] Cannot create worker ", err)
				continue
			} else {
				worker = worker2
				log.Println("[WARN-WorkerMan] New Worker .....", id)
			}
		}

		newWorkers[id] = worker
	}
	// 제거된 worker 종료하기
	for id, worker := range tempWorkers {
		worker.Stop()
		log.Println("[WARN-WorkerMan] Dispose Worker .....", id)
	}

	manager.workers = newWorkers

	for id, worker := range manager.workers {
		if !worker.IsStarted() {
			worker.Start()
			log.Println("[WARN-WorkerMan] New Worker Started .....", id)
		} else {
			log.Println("[WARN-WorkerMan] Remained Worker .....", id)
		}
	}
}
