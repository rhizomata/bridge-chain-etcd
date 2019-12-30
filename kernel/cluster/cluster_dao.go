package cluster

import (
	"fmt"
	"time"

	"github.com/rhizomata/bridge-chain-etcd/kernel/kv"
)

const (
	kvDirSys                 = "/$sys/"
	kvDirClusters            = kvDirSys + "clstrs/"
	kvPatterMemberInfo       = kvDirClusters + "%s/memb/%s"
	kvDirMemberHeartbeat     = kvDirSys + "%s/hb/"
	kvPatternMemberHeartbeat = kvDirMemberHeartbeat + "%s"
	kvPatternLeader          = kvDirClusters + "%s/leader"
)

// DAO kv store model for cluster
type DAO struct {
	cluster string
	kv      kv.KV
}

// GetLeader get leader id
func (dao *DAO) GetLeader() (leader string, err error) {
	key := fmt.Sprintf(kvPatternLeader, dao.cluster)
	bytes, err := dao.kv.GetOne(key)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// PutLeader set leader
func (dao *DAO) PutLeader(leader string) (err error) {
	key := fmt.Sprintf(kvPatternLeader, dao.cluster)
	_, err = dao.kv.Put(key, leader)
	return err
}

// GetMemberInfo ..
func (dao *DAO) GetMemberInfo(id string) (memb Member, err error) {
	key := fmt.Sprintf(kvPatterMemberInfo, dao.cluster, id)
	memb = Member{ID: id}
	// fmt.Println("********* Before GetMemberInfo:", key, memb)
	err = dao.kv.GetObject(key, &memb)
	return memb, err
}

// PutMemberInfo ..
func (dao *DAO) PutMemberInfo(memb Member) (err error) {
	key := fmt.Sprintf(kvPatterMemberInfo, dao.cluster, memb.ID)
	_, err = dao.kv.PutObject(key, memb)
	// fmt.Println("********* PutMemberInfo:", key, memb)
	return err
}

// GetHeartbeat ..
func (dao *DAO) GetHeartbeat(id string) (tm time.Time, err error) {
	bytes, err := dao.kv.GetOne(fmt.Sprintf(kvPatternMemberHeartbeat, dao.cluster, id))
	if err == nil {
		tm, err = time.Parse(time.RFC3339, string(bytes))
	}
	return tm, err
}

// GetHeartbeats  ..
func (dao *DAO) GetHeartbeats(handler func(id string, tm time.Time)) (err error) {
	dirPath := fmt.Sprintf(kvDirMemberHeartbeat, dao.cluster)

	// fmt.Println("*********** GetHeartbeats : dirPath:", dirPath)
	err = dao.kv.GetWithPrefix(dirPath,
		func(key string, value []byte) {
			a := []rune(key)
			id := string(a[len(dirPath):])

			tm, err2 := time.Parse(time.RFC3339, string(value))
			if err2 != nil {
				// log.Println("[ERROR] Parse member Heartbeat time :", err2)
			} else {
				handler(id, tm)
			}
		})
	return err
}

// PutHeartbeat ..
func (dao *DAO) PutHeartbeat(id string) (err error) {
	key := fmt.Sprintf(kvPatternMemberHeartbeat, dao.cluster, id)
	nowStr := time.Now().Format(time.RFC3339)
	_, err = dao.kv.Put(key, nowStr)
	// fmt.Println("********** PutHeartbeat:", key)
	return err
}
