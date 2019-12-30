package job

import (
	"fmt"
	"log"
)

// Organizer : Job Organizer distributes jobs to members
type Organizer interface {
	Distribute(allJobs map[string]Job, aliveMembers []string, membJobMap map[string][]string) (membJobs map[string][]string, err error)
}

type simpleOrganizer struct {
}

// NewSimpleOrganizer ..
func NewSimpleOrganizer() Organizer {
	return &simpleOrganizer{}
}

// Distribute ..
func (organizer *simpleOrganizer) Distribute(
	allJobs map[string]Job, aliveMembers []string, membJobMap map[string][]string) (membJobs map[string][]string, err error) {
	fmt.Println("======== jobOrganizer::Distribute +++++++++")

	// 1) alive하지 않은 멤버의 job 회수
	// 2) 할당되지 않은 job 식별
	// 3) 새로운 job들 alive멤버에게 할당

	unallocatedJobs := make(map[string]Job)
	for k, v := range allJobs {
		unallocatedJobs[k] = v
	}

	for _, membID := range aliveMembers {
		jobs := membJobMap[membID]
		if jobs != nil {
			for _, job := range jobs {
				delete(unallocatedJobs, job)
			}
		}
	}

	log.Println("[INFO-SimpleJobOrganizer] all jobs:", len(allJobs), ", unallocated jobs:", len(unallocatedJobs))

	avg := len(allJobs) / len(aliveMembers)
	if len(allJobs)%len(aliveMembers) > 0 {
		avg = avg + 1
	}

	newMembJobsMap := make(map[string][]string)

	// avg 보다 많은 job을 가진 멤버들 정리
	for _, membID := range aliveMembers {
		membJobs := membJobMap[membID]
		if membJobs == nil {
			membJobs = []string{}
		}

		if len(membJobs) > avg {
			jobs := membJobs[0:avg]
			remains := membJobs[avg:]
			for _, jobid := range remains {
				unallocatedJobs[jobid] = allJobs[jobid]
			}
			newMembJobsMap[membID] = jobs
		} else {
			newMembJobsMap[membID] = membJobs
		}
	}

	unallocatedJobIDs := []string{}

	for k := range unallocatedJobs {
		unallocatedJobIDs = append(unallocatedJobIDs, k)
	}

	for _, membID := range aliveMembers {
		membJobs := newMembJobsMap[membID]
		cnt := avg - len(membJobs) // cnt must be more than 0

		if cnt > 0 {
			var slice []string
			if len(unallocatedJobIDs) >= cnt {
				slice = unallocatedJobIDs[:cnt]
				unallocatedJobIDs = unallocatedJobIDs[cnt:]
			} else {
				slice = unallocatedJobIDs
				unallocatedJobIDs = []string{}
			}

			membJobs = append(membJobs, slice...)
			newMembJobsMap[membID] = membJobs
		}
		// fmt.Println("***** len(unallocatedJobIDs)=", len(unallocatedJobIDs))
		// fmt.Println("***** len(membJobs)=", len(membJobs))
	}

	// inactive한 멤버의 jobMap은 빈 어레이로 대체
	for k := range membJobMap {
		if newMembJobsMap[k] == nil {
			newMembJobsMap[k] = []string{}
		}
	}

	return newMembJobsMap, nil
}
