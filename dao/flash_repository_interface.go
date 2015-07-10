package dao

import "github.com/index0h/go-tracker/entities"

type FlashRepositoryInterface interface {
	//
	FindAll() (result []entities.Flash, err error)

	//
	FindAllByVisit(*entities.Visit) (result []entities.Flash, err error)

	//
	FindByID([16]byte) (result *entities.Flash, err error)

	//
	Insert(*entities.Flash) (err error)
}
