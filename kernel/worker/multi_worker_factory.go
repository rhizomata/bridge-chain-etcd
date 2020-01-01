package worker

import "log"

// MultiWorkerFactory  implements worker.Factory
type MultiWorkerFactory struct {
	name            string
	workerFactories map[string]Factory
}

// NewMultiWorkerFactory create MultiWorkerFactory
func NewMultiWorkerFactory(name string, factories []Factory) (factory *MultiWorkerFactory, err error) {
	factory = &MultiWorkerFactory{name: name}
	factory.workerFactories = make(map[string]Factory)

	for _, fac := range factories {
		factory.workerFactories[fac.Name()] = fac
	}

	return factory, err
}

// Name return factory.name
func (factory *MultiWorkerFactory) Name() string { return factory.name }

// NewWorker implements worker.Factory.NewWorker
func (factory *MultiWorkerFactory) NewWorker(helper *Helper) (wroker Worker, err error) {
	multiWorker := &MultiWorker{id: helper.ID()}
	multiWorker.workers = make(map[string]Worker)

	for name, fac := range factory.workerFactories {
		subHelper := helper.CreateChildHelper(factory.name+"-"+name, helper.job)
		subWorker, err := fac.NewWorker(subHelper)
		if err != nil {
			log.Println("[ERROR-MW-", factory.Name(), "] Create Sub Worker ", name, err)
		}
		multiWorker.addWorker(subWorker)
	}

	return multiWorker, nil
}

// MultiWorker imeplements Worker, which has many sub workers
type MultiWorker struct {
	id      string
	started bool
	workers map[string]Worker
}

func (multiWorker *MultiWorker) addWorker(worker Worker) {
	multiWorker.workers[worker.ID()] = worker
}

// ID return multiWorker.id
func (multiWorker *MultiWorker) ID() string { return multiWorker.id }

// IsStarted return multiWorker.started
func (multiWorker *MultiWorker) IsStarted() bool { return multiWorker.started }

// Start ..
func (multiWorker *MultiWorker) Start() (err error) {
	for id, worker := range multiWorker.workers {
		err = worker.Start()
		if err != nil {
			log.Println("[ERROR-MultiWorker-", multiWorker.ID(), "] Start Worker ", id, err)
			multiWorker.Stop()
			break
		}
	}
	if err == nil {
		multiWorker.started = true
		log.Println("[INFO-MultiWorker-", multiWorker.ID(), "] Started")
	}
	return err
}

// Stop ..
func (multiWorker *MultiWorker) Stop() error {
	for _, worker := range multiWorker.workers {
		if worker.IsStarted() {
			worker.Stop()
		}
	}
	log.Println("[INFO-MultiWorker-", multiWorker.ID(), "] Stoped.")
	return nil
}
