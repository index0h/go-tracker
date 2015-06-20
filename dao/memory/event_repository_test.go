package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/index0h/go-tracker/entities"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&Repository{})
}

func Test_EventRepository_NewEventRepository(t *testing.T) {
	nested := new(nestedEventRepository)
	result := []*entities.Event{}
	nested.On("FindAll").Return(result, nil)

	repository, err := NewEventRepository(nested)

	assert.NotNil(t, repository)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func Test_EventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewEventRepository(nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

type nestedEventRepository struct {
	mock.Mock
}

func (repository *nestedEventRepository) FindAll() (result []*entities.Event, err error) {
	args := repository.Called()

	raw := args.Get(0)
	result, _ = raw.([]*entities.Event)

	return result, args.Error(1)
}

func (repository *nestedEventRepository) FindAllByVisit(visit *entities.Visit) (result []*entities.Event, err error) {
	args := repository.Called(visit)

	raw := args.Get(0)
	result, _ = raw.([]*entities.Event)

	return result, args.Error(1)
}

func (repository *nestedEventRepository) FindByID(eventID [16]byte) (result *entities.Event, err error) {
	args := repository.Called(eventID)

	raw := args.Get(0)
	result, _ = raw.(*entities.Event)

	return result, args.Error(1)
}

func (repository *nestedEventRepository) Insert(event *entities.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}

func (repository *nestedEventRepository) Update(event *entities.Event) (err error) {
	args := repository.Called(event)

	return args.Error(0)
}
