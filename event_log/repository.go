package event

import (
	eventLogEntities "github.com/index0h/go-tracker/event_log/entities"
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

type Repository interface {
	//
	FindAll() (result []eventLogEntities.EventLog, err error)

	//
	FindAllByVisit(*visitEntities.Visit) (result []eventLogEntities.EventLog, err error)

	//
	FindByID(uuid.UUID) (result *eventLogEntities.EventLog, err error)

	//
	Insert(*eventLogEntities.EventLog) (err error)

	//
	Update(*eventLogEntities.EventLog) (err error)
}
