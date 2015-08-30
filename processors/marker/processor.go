package marker

import (
	eventEntity "github.com/index0h/go-tracker/modules/event/entity"
	flashEntity "github.com/index0h/go-tracker/modules/flash/entity"
	visitEntity "github.com/index0h/go-tracker/modules/visit/entity"
)

type Dummy struct {
	priority int
}

func New(priority int) *Dummy {
	return &Dummy{priority: priority}
}

func (processor *Dummy) GetPriority() int {
	return processor.priority
}

func (processor *Dummy) Process(
	flash *flashEntity.Flash,
	event *eventEntity.Event,
	visit *visitEntity.Visit,
) *flashEntity.Flash {
	return flash
}

// marker.set.registrationData = firstLogin
// marker.push.transactionIds = transaction_id
