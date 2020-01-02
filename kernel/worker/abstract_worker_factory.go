package worker

import "errors"

// AbstractWorkerFactory implements worker.Factory, job data format : #factoryName:data
type AbstractWorkerFactory struct {
	name            string
	workerFactories map[string]Factory
}

// Name return factory.name
func (abstractFactory *AbstractWorkerFactory) Name() string { return abstractFactory.name }

// AddFactory add worker factory
func (abstractFactory *AbstractWorkerFactory) AddFactory(factory Factory) {
	abstractFactory.workerFactories[factory.Name()] = factory
}

// GetFactory get worker factory
func (abstractFactory *AbstractWorkerFactory) GetFactory(name string) (factory Factory, err error) {
	factory = abstractFactory.workerFactories[name]
	if factory == nil {
		err = errors.New("Factory not found for " + name)
	}
	return factory, err
}

// NewAbstractWorkerFactory create AbstractWorkerFactory
func NewAbstractWorkerFactory(name string) (factory *AbstractWorkerFactory) {
	factory = &AbstractWorkerFactory{name: name}
	factory.workerFactories = make(map[string]Factory)
	return factory
}

const (
	sharp = byte('#')
	colon = byte(':')
)

// NewWorker implements worker.Factory.NewWorker
func (abstractFactory *AbstractWorkerFactory) NewWorker(helper *Helper) (wroker Worker, err error) {
	jobData := helper.Job()

	if jobData[0] == sharp {
		index := -1
		for i, b := range jobData {
			index = i
			if b == colon {
				break
			}
		}
		if index < 2 {
			err = errors.New("Job Data must be started with '#factory-name:'")
			return nil, err
		}
		factoryName := string(jobData[1:index])

		factory, err := abstractFactory.GetFactory(factoryName)
		if err != nil {
			return nil, err
		}
		helper.job = jobData[index+1:]
		return factory.NewWorker(helper)
	}
	err = errors.New("Job Data must be started with '#factory-name:'")
	return nil, err
}
