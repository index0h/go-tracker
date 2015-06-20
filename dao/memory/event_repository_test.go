package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/stretchr/testify/assert"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&Repository{})
}

func Test_EventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewEventRepository(nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}
