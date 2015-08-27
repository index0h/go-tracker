package memory

import (
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (repository *MockEventRepository) FindAll(limit int64, offset int64) (result []*entities.Event, err error) {
	args := repository.Called(limit, offset)

	raw := args.Get(0)
	result, _ = raw.([]*entities.Event)

	return result, args.Error(1)
}

func (repository *MockEventRepository) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	args := repository.Called(visit)

	raw := args.Get(0)
	result, _ = raw.([]*entities.Event)

	return result, args.Error(1)
}

func (repository *MockEventRepository) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	args := repository.Called(eventID)

	raw := args.Get(0)
	result, _ = raw.(*entities.Event)

	return result, args.Error(1)
}

func (repository *MockEventRepository) Insert(event *entities.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}

func (repository *MockEventRepository) Update(event *entities.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}
