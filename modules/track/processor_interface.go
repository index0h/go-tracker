package track

import (
	event "github.com/index0h/go-tracker/modules/event/entity"
	flash "github.com/index0h/go-tracker/modules/flash/entity"
	visit "github.com/index0h/go-tracker/modules/visit/entity"
)

type ProcessorInterface interface {
	//
	Process(*flash.Flash, *event.Event, *visit.Visit) *flash.Flash

	GetPriority() int
}
