package dao

import "github.com/index0h/go-tracker/entities"

type MarkRepositoryInterface interface {

	//
	FindByID([16]byte) (result *entities.Mark, err error)

	//
	FindByClientID(string) (result *entities.Mark, err error)

	//
	FindAll(limit int64, offset int64) (result []*entities.Mark, err error)

	//
	Insert(*entities.Mark) (err error)

	//
	Update(*entities.Mark) (err error)
}
