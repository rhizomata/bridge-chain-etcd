package job

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Job job data structure
type Job struct {
	ID   string
	Data []byte
}

// NewJob ..
func NewJob(data []byte) Job {
	uuid := uuid.New()
	return Job{ID: uuid.String(), Data: data}
}

// GetAsString Get data as string
func (job *Job) GetAsString() string {
	return string(job.Data)
}

// GetAsObject Get data as interface
func (job *Job) GetAsObject(obj interface{}) error {
	err := json.Unmarshal(job.Data, obj)
	return err
}
