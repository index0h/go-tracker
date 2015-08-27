package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&EventRepository{})
}

func TestEventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewEventRepository(nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}
