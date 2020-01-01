package worker 

import (
	"fmt"
	"log"

	"github.com/rhizomata/bridge-chain-etcd/kernel/kv"
)

const (
	kvDirSys            = "/$sys/"
	kvDirClusters       = kvDirSys + "clstrs/"
	kvPatternCheckpoint = kvDirClusters + "%s/checkpoint/%s"
	kvPatternDataJobID  = kvDirClusters + "%s/data/%s/"
	kvPatternData       = kvPatternDataJobID + "%s"
)

// DAO kv store model for cluster
type DAO struct {
	cluster string
	kv      kv.KV
}

// PutCheckpoint ..
func (dao *DAO) PutCheckpoint(jobid string, checkpoint interface{}) error {
	_, err := dao.kv.PutObject(fmt.Sprintf(kvPatternCheckpoint, dao.cluster, jobid), checkpoint)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutCheckpoint", err)
	}
	return err
}

// GetCheckpoint ..
func (dao *DAO) GetCheckpoint(jobid string, checkpoint interface{}) error {
	err := dao.kv.GetObject(fmt.Sprintf(kvPatternCheckpoint, dao.cluster, jobid), checkpoint)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetCheckpoint ", err)
	}
	return err
}

// PutData ..
func (dao *DAO) PutData(jobid string, rowID string, data interface{}) error {
	_, err := dao.kv.PutObject(fmt.Sprintf(kvPatternData, dao.cluster, jobid, rowID), data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] PutData", err)
	}
	return err
}

// GetData ..
func (dao *DAO) GetData(jobid string, rowID string, data interface{}) error {
	err := dao.kv.GetObject(fmt.Sprintf(kvPatternData, dao.cluster, jobid, rowID), data)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetData ", err)
	}
	return err
}

// DeleteData ..
func (dao *DAO) DeleteData(jobid string, rowID string) error {
	_, err := dao.kv.DeleteOne(fmt.Sprintf(kvPatternData, dao.cluster, jobid, rowID))
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetData ", err)
	}
	return err
}

// GetDataWithJobID ..
func (dao *DAO) GetDataWithJobID(jobid string, handler func(key string, value []byte)) error {
	err := dao.kv.GetWithPrefix(fmt.Sprintf(kvPatternDataJobID, dao.cluster, jobid), handler)
	if err != nil {
		log.Println("[ERROR-WorkerDao] GetDataWithJobID ", err)
	}
	return err
}
