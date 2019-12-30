package cluster

import (
	"fmt"
	"log"
	"time"

	"github.com/rhizomata/bridge-chain-etcd/kernel/kv"
	"github.com/rhizomata/bridge-chain-etcd/kernel/model"
)

// Manager cluster manager
type Manager struct {
	cluster              *Cluster
	dao                  *DAO
	config               model.Config
	running              bool
	memberChangeHandler  func(aliveMembers []string)
	healthCheckDelegator func(memb *Member) bool
}

// NewManager create cluster
func NewManager(localid string, config model.Config, kv kv.KV) *Manager {
	cluster := newCluster(config.Cluster)
	dao := DAO{cluster: config.Cluster, kv: kv}

	manager := new(Manager)
	manager.cluster = cluster
	manager.dao = &dao
	manager.config = config

	localMemb := Member{Cluster: cluster.name, ID: localid, Name: config.Name, DaemonURL: config.GetDaemonURL()}
	localMemb.setLocal(true)
	localMemb.setAlive(true)

	cluster.putMember(&localMemb)

	cluster.localMember = &localMemb

	return manager
}

// SetMemberChangeHandler set memberChangeHandler
// 이 데몬이 리더 자격을 얻었을 때 사용할 Job Organize용 핸들러
func (manager *Manager) SetMemberChangeHandler(memberChangeHandler func(aliveMembers []string)) {
	manager.memberChangeHandler = memberChangeHandler
}

// SetHealthCheckDelegator ..
func (manager *Manager) SetHealthCheckDelegator(healthCheckDelegator func(memb *Member) bool) {
	manager.healthCheckDelegator = healthCheckDelegator
}

// Start start goroutins
func (manager *Manager) Start() {
	manager.running = true

	err := manager.dao.PutMemberInfo(*manager.cluster.localMember)
	if err != nil {
		log.Fatal("[FATAL] Cannot send PutMemberInfo.", err)
	}
	err = manager.dao.PutHeartbeat(manager.cluster.localMember.ID)

	if err != nil {
		log.Fatal("[FATAL] Cannot send heartbeat.", err)
	}

	go func() {
		for manager.running {
			time.Sleep(time.Duration(manager.config.HeartbeatInterval))
			err := manager.dao.PutHeartbeat(manager.cluster.localMember.ID)
			if err != nil {
				log.Fatal("[FATAL] Cannot send heartbeat.", err)
			}
		}
	}()

	go func() {
		for manager.running {
			err := manager.dao.GetHeartbeats(manager.handleHeartbeat)
			if err != nil {
				log.Fatal("[FATAL] Cannot check heartbeats.", err)
			}
			manager.checkLeader()
			time.Sleep(time.Duration(manager.config.CheckHeartbeatInterval))
		}
	}()

	log.Println("[INFO-Cluster] Start Cluster Manager.")
}

// Dispose stop goroutins
func (manager *Manager) Dispose() {
	manager.running = false
	log.Println("[WARN-Cluster] Dispose Cluster Manager.")
}

func (manager *Manager) handleHeartbeat(id string, tm time.Time) {
	if manager.cluster.Local().ID == id {
		manager.cluster.Local().setHeartBeat(tm)
	}
	memb := manager.cluster.GetMember(id)

	changed := false

	if memb == nil {
		memb2, err := manager.dao.GetMemberInfo(id)
		if err != nil {
			log.Println("[ERROR-Cluster] Cannot find member info ", id, err)
		}
		memb = &memb2
		manager.cluster.putMember(memb)
		changed = true
	}

	alive := false

	if memb.IsLocal() {
		alive = true
	} else {
		if memb.HeartBeat().IsZero() || memb.HeartBeat() == tm {
			now := time.Now()
			duration := now.Sub(tm)
			if manager.healthCheckDelegator != nil {
				alive = manager.healthCheckDelegator(memb)
			} else {
				alive = duration.Seconds() < float64(manager.config.AliveThreasholdSeconds)
			}
		} else {
			alive = true
		}
	}

	memb.setHeartBeat(tm)

	if memb.IsAlive() != alive {
		changed = true
		memb.setAlive(alive)
	}

	// fmt.Println("***** handleHeartbeat :: ", memb)

	if changed {
		manager.onMemberChanged(memb)
	}
}

func (manager *Manager) checkLeader() {
	leaderID, err := manager.dao.GetLeader()
	if err != nil {
		log.Println("[ERROR] Get Leader ", err)
	}

	oldLeader := manager.cluster.leader

	if oldLeader != nil {
		if oldLeader.ID == leaderID {
			if oldLeader.IsAlive() {
				return
			}
		} else {
			oldLeader.setLeader(false)
			manager.cluster.leader = nil
		}
	}

	var leader *Member

	if leaderID != "" {
		leader = manager.cluster.GetMember(leaderID)
		if leader == nil {
			leader = manager.electLeader()
		} else if !leader.IsAlive() {
			leader.setLeader(false)
			leader = manager.electLeader()
		}
	} else {
		leader = manager.electLeader()
	}

	leader.setLeader(true)
	manager.cluster.leader = leader

	manager.onLeaderChanged(leader)
}

func (manager *Manager) electLeader() *Member {
	members := manager.cluster.GetSortedMembers()

	fmt.Println("****** electLeader:: len(members) ", len(members))

	for _, id := range members {
		memb := manager.cluster.GetMember(id)
		fmt.Println("    ****** electLeader:: member ", id, memb)
		if memb.IsAlive() {
			manager.dao.PutLeader(id)
			return memb
		}
	}

	local := manager.cluster.Local()
	manager.dao.PutLeader(local.ID)
	return local
}

// IsLeader : returns whether this kernel is leader.
func (manager *Manager) IsLeader() bool {
	return manager.cluster.localMember.IsLeader()
}

func (manager *Manager) onMemberChanged(memb *Member) {
	if manager.cluster.localMember.IsLeader() && manager.memberChangeHandler != nil {
		log.Println("[INFO-Cluster] Member Changed::", memb)
		manager.memberChanged()
	}
}

func (manager *Manager) onLeaderChanged(leader *Member) {
	if manager.cluster.localMember.IsLeader() && manager.memberChangeHandler != nil {
		log.Println("[INFO-Cluster] Leader changed. I'm the leader")
		manager.memberChanged()
	}
}

func (manager *Manager) memberChanged() {
	manager.memberChangeHandler(manager.cluster.GetAliveMemberIDs())
}

// GetCluster get cluster
func (manager *Manager) GetCluster() *Cluster {
	return manager.cluster
}
