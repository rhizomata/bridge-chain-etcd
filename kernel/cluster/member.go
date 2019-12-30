package cluster

import "time"

// Member member info
type Member struct {
	Cluster   string `json:"cluster"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	DaemonURL string `json:"url"`
	heartbeat time.Time
	leader    bool
	alive     bool
	local     bool
}

//HeartBeat return member's last heartbeat time
func (memb *Member) HeartBeat() time.Time {
	return memb.heartbeat
}

//setHeartBeat set member's last heartbeat time
func (memb *Member) setHeartBeat(time time.Time) {
	memb.heartbeat = time
}

//IsLeader return whether member is leader
func (memb *Member) IsLeader() bool {
	return memb.leader
}

//setLeader Set member as leader
func (memb *Member) setLeader(leader bool) {
	memb.leader = leader
}

//IsAlive return whether member is alive
func (memb *Member) IsAlive() bool {
	return memb.alive
}

//setAlive Set member alive
func (memb *Member) setAlive(alive bool) {
	memb.alive = alive
}

//IsLocal return whether member is alive
func (memb *Member) IsLocal() bool {
	return memb.local
}

//setLocal Set member alive
func (memb *Member) setLocal(local bool) {
	memb.local = local
}
