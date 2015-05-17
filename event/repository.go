package event

import (
	eventEntities "github.com/index0h/go-tracker/event/entities"
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type Repository interface {
	//
	FindAll() (result []eventEntities.Event, err error)

	//
	FindAllByVisit(*visitEntities.Visit) (result []eventEntities.Event, err error)

	//
	FindByID(uuid.UUID) (result *eventEntities.Event, err error)

	//
	Insert(*eventEntities.Event) (err error)

	//
	Update(*eventEntities.Event) (err error)
}
