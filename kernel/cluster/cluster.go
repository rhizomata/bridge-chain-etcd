package cluster

import (
	"sort"
)

// Cluster ..
type Cluster struct {
	name        string
	membIDs     []string
	members     map[string]*Member
	localMember *Member
	leader      *Member
}

func newCluster(name string) *Cluster {
	cluster := Cluster{name: name}
	cluster.membIDs = []string{}
	cluster.members = make(map[string]*Member)

	return &cluster
}

// Name get name
func (cluster *Cluster) Name() string {
	return cluster.name
}

func (cluster *Cluster) putMember(memb *Member) {
	if _, ok := cluster.members[memb.ID]; !ok {
		cluster.membIDs = append(cluster.membIDs, memb.ID)
		sort.Strings(cluster.membIDs)
	}
	// fmt.Println("********* putMember :: ", cluster.members, memb.ID, memb)
	cluster.members[memb.ID] = memb
}

func (cluster *Cluster) removeMember(id string) {
	index := -1
	for i, mid := range cluster.membIDs {
		if mid == id {
			index = i
			break
		}
	}

	if index > -1 {
		if index <= len(cluster.membIDs)-1 {
			copy(cluster.membIDs[index:], cluster.membIDs[index+1:])
		}
		cluster.membIDs = cluster.membIDs[0 : len(cluster.membIDs)-1]
	}

	delete(cluster.members, id)
}

// GetMember get member with given name
func (cluster *Cluster) GetMember(id string) *Member {
	memb := cluster.members[id]
	return memb
}

// GetSortedMembers get all member ids
func (cluster *Cluster) GetSortedMembers() []string {
	return cluster.membIDs
}

// GetAliveMembers get active members
func (cluster *Cluster) GetAliveMembers() []*Member {
	membs := []*Member{}
	for _, memb := range cluster.members {
		if memb.IsAlive() {
			membs = append(membs, memb)
		}
	}
	return membs
}

// GetAliveMemberIDs get active member IDs
func (cluster *Cluster) GetAliveMemberIDs() []string {
	membs := []string{}
	for id, memb := range cluster.members {
		if memb.IsAlive() {
			membs = append(membs, id)
		}
	}
	return membs
}

// Leader get Leader
func (cluster *Cluster) Leader() *Member {
	return cluster.leader
}

// Local get localMember
func (cluster *Cluster) Local() *Member {
	return cluster.localMember
}
