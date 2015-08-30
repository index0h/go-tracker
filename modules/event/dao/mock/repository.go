package mock

import (
	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Event, err error) {
	args := repository.Called(limit, offset)

	raw := args.Get(0)
	result, _ = raw.([]*entity.Event)

	return result, args.Error(1)
}

func (repository *Repository) FindAllByFields(data types.Hash) (result []*entity.Event, err error) {
	args := repository.Called(data)

	raw := args.Get(0)
	result, _ = raw.([]*entity.Event)

	return result, args.Error(1)
}

func (repository *Repository) FindByID(eventID [16]byte) (result *entity.Event, err error) {
	args := repository.Called(eventID)

	raw := args.Get(0)
	result, _ = raw.(*entity.Event)

	return result, args.Error(1)
}

func (repository *Repository) Insert(event *entity.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}

func (repository *Repository) Update(event *entity.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}
