package dao

import "github.com/index0h/go-tracker/entities"

type EventLogRepositoryInterface interface {
	//
	FindAll() (result []entities.EventLog, err error)

	//
	FindAllByVisit(*entities.Visit) (result []entities.EventLog, err error)

	//
	FindByID([16]byte) (result *entities.EventLog, err error)

	//
	Insert(*entities.EventLog) (err error)

	//
	Update(*entities.EventLog) (err error)
}
