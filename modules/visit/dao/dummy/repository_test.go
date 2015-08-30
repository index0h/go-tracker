package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Interface(t *testing.T) {
	func(event visit.RepositoryInterface) {}(&Repository{})
}

func TestRepository_Verify(t *testing.T) {
	checkRepository := Repository{}

	ok, err := checkRepository.Verify(uuid.New().Generate(), "12345")

	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestRepository_Verify_EmptySessionID(t *testing.T) {
	checkRepository := Repository{}

	ok, err := checkRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Verify_EmptyClientID(t *testing.T) {
	checkRepository := Repository{}

	ok, err := checkRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Insert(t *testing.T) {
	checkRepository := Repository{}

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fileds := types.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entity.NewVisit(visitID, timestamp, sessionID, clientID, fileds)

	err := checkRepository.Insert(visit)

	assert.Nil(t, err)
}

func TestRepository_Insert_Nil(t *testing.T) {
	checkRepository := Repository{}

	err := checkRepository.Insert(nil)

	assert.NotNil(t, err)
}
