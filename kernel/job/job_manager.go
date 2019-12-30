package job

import (
	"log"

	"github.com/rhizomata/bridge-chain-etcd/kernel/kv"
)

// Manager manager for jobs
type Manager struct {
	cluster             string
	localid             string
	dao                 *DAO
	jobWatchHandler     func(job *Job)
	jobWatcher          *kv.Watcher
	membJobWatchHandler func(jobids []string)
	membJobWatcher      *kv.Watcher
}

// NewManager ..
func NewManager(cluster string, localid string, kv kv.KV) *Manager {
	manager := Manager{cluster: cluster, localid: localid, dao: &DAO{cluster: cluster, kv: kv}}
	return &manager
}

// SetMembJobWatchHandler : Set JobOrganizer
func (manager *Manager) SetMembJobWatchHandler(handler func(jobids []string)) {
	manager.membJobWatchHandler = handler
}

// SetJobWatchHandler : Set JobOrganizer
func (manager *Manager) SetJobWatchHandler(handler func(job *Job)) {
	manager.jobWatchHandler = handler
}

// Start watchers ..
func (manager *Manager) Start() {
	manager.jobWatcher = manager.dao.WatchJobs(
		func(jobid string, data []byte) {
			if manager.jobWatchHandler != nil {
				job := Job{ID: jobid, Data: data}
				manager.jobWatchHandler(&job)
			}
		})

	manager.membJobWatcher = manager.dao.WatchMemberJobs(manager.localid,
		func(jobids []string) {
			if manager.membJobWatchHandler != nil {
				manager.membJobWatchHandler(jobids)
			}
		})

}

// Dispose watchers ..
func (manager *Manager) Dispose() {
	manager.jobWatcher.Stop()
	manager.membJobWatcher.Stop()
}

// AddJob ..
func (manager *Manager) AddJob(job Job) error {
	return manager.dao.PutJob(job.ID, job.Data)
}

// RemoveJob ..
func (manager *Manager) RemoveJob(jobID string) error {
	return manager.dao.RemoveJob(jobID)
}

// GetJob ..
func (manager *Manager) GetJob(jobID string) (job Job, err error) {
	return manager.dao.GetJob(jobID)
}

// GetAllJobIDs ..
func (manager *Manager) GetAllJobIDs() (jobIDs []string, err error) {
	return manager.dao.GetAllJobIDs()
}

// GetAllJobs ..
func (manager *Manager) GetAllJobs() (jobs map[string]Job, err error) {
	return manager.dao.GetAllJobs()
}

// GetMemberJobIDs ..
func (manager *Manager) GetMemberJobIDs(membID string) (jobIDs []string, err error) {
	return manager.dao.GetMemberJobs(membID)
}

// GetAllMemberJobIDs : returns member-JobIDs Map
func (manager *Manager) GetAllMemberJobIDs() (membJobMap map[string][]string, err error) {
	return manager.dao.GetAllMemberJobIDs()
}

// SetMemberJobIDs ..
func (manager *Manager) SetMemberJobIDs(membID string, jobIDs []string) (err error) {
	return manager.dao.PutMemberJobs(membID, jobIDs)
}

// GetMemberJobs ..
func (manager *Manager) GetMemberJobs(membID string) (jobs []Job, err error) {
	jobIDs, err := manager.dao.GetMemberJobs(membID)
	if err != nil {
		log.Println("[ERROR] Cannot retrieve member jobs", err)
		return []Job{}, err
	}
	jobs = []Job{}
	for _, jobID := range jobIDs {
		job, err2 := manager.dao.GetJob(jobID)
		if err2 == nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}
