package memory

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_Interface(t *testing.T) {
	func(event event.RepositoryInterface) {}(&Repository{})
}

func TestEventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewRepository(nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}
