package dao

import "github.com/index0h/go-tracker/entities"

type FlashRepositoryInterface interface {
	//
	FindByID([16]byte) (result *entities.Flash, err error)

	//
	FindAll(limit int64, offset int64) (result []*entities.Flash, err error)

	//
	FindAllByVisitID(visitID [16]byte) (result []*entities.Flash, err error)

	//
	FindAllByEventID(eventID [16]byte, limit int64, offset int64) (result []*entities.Flash, err error)

	//
	Insert(*entities.Flash) (err error)
}
