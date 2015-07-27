package dummy

import "github.com/index0h/go-tracker/entities"

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
	flash *entities.Flash,
	event *entities.Event,
	visit *entities.Visit,
) *entities.Flash {
	return flash
}
