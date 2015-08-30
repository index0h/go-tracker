package mark

import (
	"github.com/index0h/go-tracker/modules/mark/entity"
	"github.com/index0h/go-tracker/share/types"
)

type RepositoryInterface interface {

	//
	FindByID(types.UUID) (result *entity.Mark, err error)

	//
	FindByClientID(string) (result *entity.Mark, err error)

	//
	FindBySessionID(types.UUID) (result *entity.Mark, err error)

	//
	FindAll(limit int64, offset int64) (result []*entity.Mark, err error)

	//
	Insert(*entity.Mark) (err error)

	//
	Update(*entity.Mark) (err error)
}
