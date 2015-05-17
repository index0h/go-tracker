package event

import (
	"github.com/index0h/go-tracker/uuid"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	eventLogEntities "github.com/index0h/go-tracker/eventLog/entities"
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