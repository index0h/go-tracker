package event

import (
	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
)

type RepositoryInterface interface {
	//
	FindAll(limit int64, offset int64) (result []*entity.Event, err error)

	//
	FindAllByFields(fields types.Hash) (result []*entity.Event, err error)

	//
	FindByID(types.UUID) (result *entity.Event, err error)

	//
	Insert(*entity.Event) (err error)

	//
	Update(*entity.Event) (err error)
}
