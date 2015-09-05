package memory

import (
	"testing"

	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/dao/mock"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Interface(t *testing.T) {
	func(visit visit.RepositoryInterface) {}(&Repository{})
}

func TestRepository_NewRepository_EmptyClient(t *testing.T) {
	repository, err := NewRepository(nil, 10)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestRepository_Verify(t *testing.T) {
	nested := new(mock.Repository)
	checkRepository, _ := NewRepository(nested, 10)

	sessionID := uuid.New().Generate()
	clientID := "12345"

	nested.On("Verify", sessionID, clientID).Return(true, nil)

	ok, err := checkRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestRepository_Verify_Cached(t *testing.T) {
	nested := new(mock.Repository)
	checkRepository, _ := NewRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := types.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entity.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	checkRepository.Insert(visit)

	ok, err := checkRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestRepository_Verify_RegisteredByAnother(t *testing.T) {
	nested := new(mock.Repository)
	checkRepository, _ := NewRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := types.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entity.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	checkRepository.Insert(visit)

	ok, err := checkRepository.Verify(sessionID, "AnotherClientID")

	assert.False(t, ok)
	assert.NotNil(t, err)
	nested.AssertExpectations(t)
}

func TestRepository_Verify_EmptySessionID(t *testing.T) {
	checkRepository, _ := NewRepository(new(mock.Repository), 10)

	ok, err := checkRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Verify_EmptyClientID(t *testing.T) {
	checkRepository, _ := NewRepository(new(mock.Repository), 10)

	ok, err := checkRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestRepository_Insert(t *testing.T) {
	nested := new(mock.Repository)
	checkRepository, _ := NewRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := types.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entity.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	err := checkRepository.Insert(visit)

	_, okClientID := checkRepository.cache.Get(sessionID)

	assert.True(t, okClientID)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestRepository_Insert_Nil(t *testing.T) {
	nested := new(mock.Repository)
	checkRepository, _ := NewRepository(nested, 10)

	err := checkRepository.Insert(nil)

	assert.NotNil(t, err)
}
