package components

import (
	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/entities"
)

type EventLogRepositoryInterface interface {
	//
	FindAll() (result []entities.EventLog, err error)

	//
	FindAllByVisit(*entities.Visit) (result []entities.EventLog, err error)

	//
	FindByID(common.UUID) (result *entities.EventLog, err error)

	//
	Insert(*entities.EventLog) (err error)

	//
	Update(*entities.EventLog) (err error)
}
