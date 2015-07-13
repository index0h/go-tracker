package processors

import "github.com/index0h/go-tracker/entities"

type DummyProcessor struct {
	priority int
}

func NewDummy(priority int) *DummyProcessor {
	return &DummyProcessor{priority: priority}
}

func (processor *DummyProcessor) Process(
	flash *entities.Flash,
	event *entities.Event,
	visit *entities.Visit,
) *entities.Flash {
	return flash
}

func (processor *DummyProcessor) GetPriority() int {
	return processor.priority
}
