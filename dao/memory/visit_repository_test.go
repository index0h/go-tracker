package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func TestVisitRepository_Interface(t *testing.T) {
	func(visit dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func TestVisitRepository_NewVisitRepository_EmptyClient(t *testing.T) {
	repository, err := NewVisitRepository(nil, 10)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestVisitRepository_Verify(t *testing.T) {
	nested := new(MockVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	sessionID := uuid.New().Generate()
	clientID := "12345"

	nested.On("Verify", sessionID, clientID).Return(true, nil)

	ok, err := checkVisitRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_Cached(t *testing.T) {
	nested := new(MockVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fileds := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fileds)

	nested.On("Insert", visit).Return(nil)

	checkVisitRepository.Insert(visit)

	ok, err := checkVisitRepository.Verify(sessionID, clientID)

	assert.True(t, ok)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_RegisteredByAnother(t *testing.T) {
	nested := new(MockVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	checkVisitRepository.Insert(visit)

	ok, err := checkVisitRepository.Verify(sessionID, "AnotherClientID")

	assert.False(t, ok)
	assert.NotNil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Verify_EmptySessionID(t *testing.T) {
	checkVisitRepository, _ := NewVisitRepository(new(MockVisitRepository), 10)

	ok, err := checkVisitRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVisitRepository_Verify_EmptyClientID(t *testing.T) {
	checkVisitRepository, _ := NewVisitRepository(new(MockVisitRepository), 10)

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVisitRepository_Insert(t *testing.T) {
	nested := new(MockVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := entities.Hash{"data": "here"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, fields)

	nested.On("Insert", visit).Return(nil)

	err := checkVisitRepository.Insert(visit)

	_, okClientID := checkVisitRepository.sessionToClient.Get(sessionID)
	_, okSessionID := checkVisitRepository.clientToSession.Get(clientID)

	assert.True(t, okClientID)
	assert.True(t, okSessionID)
	assert.Nil(t, err)
	nested.AssertExpectations(t)
}

func TestVisitRepository_Insert_Nil(t *testing.T) {
	nested := new(MockVisitRepository)
	checkVisitRepository, _ := NewVisitRepository(nested, 10)

	err := checkVisitRepository.Insert(nil)

	assert.NotNil(t, err)
}
