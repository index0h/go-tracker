package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&EventRepository{})
}

func TestEventRepository_NewEventRepository(t *testing.T) {
	nested := new(MockEventRepository)
	result := []*entities.Event{}
	nested.On("FindAll").Return(result, nil)

	repository, err := NewEventRepository(nested)

	assert.NotNil(t, repository)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestEventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewEventRepository(nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}
