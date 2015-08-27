package dao

import "github.com/index0h/go-tracker/entities"

type ProcessorInterface interface {
	//
	Process(*entities.Flash, *entities.Event, *entities.Visit) *entities.Flash

	GetPriority() int
}
