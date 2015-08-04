package marker

import "github.com/index0h/go-tracker/entities"

type Marker struct {
	priority int
}

func New(priority int) *Marker {
	return &Marker{priority: priority}
}

func (processor *Marker) GetPriority() int {
	return processor.priority
}

func (processor *Marker) Process(
	flash *entities.Flash,
	event *entities.Event,
	visit *entities.Visit,
) *entities.Flash {
	return flash
}

// marker.set.registrationData = firstLogin
// marker.push.transactionIds = transaction_id
