package dao

import "github.com/index0h/go-tracker/entities"

type EventRepositoryInterface interface {
	//
	FindAll() (result []*entities.Event, err error)

	//
	FindAllByVisit(*entities.Visit) (result []*entities.Event, err error)

	//
	FindByID([16]byte) (result *entities.Event, err error)

	//
	Insert(*entities.Event) (err error)

	//
	Update(*entities.Event) (err error)
}
