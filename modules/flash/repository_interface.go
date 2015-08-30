package flash

import (
	"github.com/index0h/go-tracker/modules/flash/entity"
	"github.com/index0h/go-tracker/share/types"
)

type RepositoryInterface interface {
	//
	FindByID(types.UUID) (result *entity.Flash, err error)

	//
	FindAll(limit int64, offset int64) (result []*entity.Flash, err error)

	//
	FindAllByVisitID(visitID types.UUID) (result []*entity.Flash, err error)

	//
	FindAllByEventID(eventID types.UUID, limit int64, offset int64) (result []*entity.Flash, err error)

	//
	Insert(*entity.Flash) (err error)
}
