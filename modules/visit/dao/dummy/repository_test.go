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
	func(visit.RepositoryInterface) {}(&Repository{})
}

func TestRepository_FindByID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindByID(uuid.New().Generate())

	assert.Nil(t, found)
	assert.Nil(t, err)
}

func TestRepository_FindByID_EmptyVisitID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindByID(types.UUID{})

	assert.Nil(t, found)
	assert.NotNil(t, err)
}

func TestRepository_FindAll(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindAll(0, 0)

	assert.Empty(t, found)
	assert.Nil(t, err)
}

func TestRepository_FindAllBySessionID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindAllBySessionID(uuid.New().Generate(), 0, 0)

	assert.Empty(t, found)
	assert.Nil(t, err)
}

func TestRepository_FindAllBySessionID_EmptySessionID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindAllBySessionID(types.UUID{}, 0, 0)

	assert.Empty(t, found)
	assert.NotNil(t, err)
}

func TestRepository_FindAllByClientID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindAllByClientID("some client id", 0, 0)

	assert.Empty(t, found)
	assert.Nil(t, err)
}

func TestRepository_FindAllByClientID_EmptyClientID(t *testing.T) {
	checkRepository := NewRepository()

	found, err := checkRepository.FindAllByClientID("", 0, 0)

	assert.Empty(t, found)
	assert.NotNil(t, err)
}

func TestRepository_Verify(t *testing.T) {
	checkRepository := NewRepository()

	ok, err := checkRepository.Verify(uuid.New().Generate(), "12345")

	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestRepository_Verify_EmptySessionID(t *testing.T) {
	checkRepository := NewRepository()

	ok, err := checkRepository.Verify(types.UUID{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Verify_EmptyClientID(t *testing.T) {
	checkRepository := NewRepository()

	ok, err := checkRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Insert(t *testing.T) {
	checkRepository := NewRepository()

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
	checkRepository := NewRepository()

	err := checkRepository.Insert(nil)

	assert.NotNil(t, err)
}
