package memory

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_VisitRepository_FindClientID_Empty(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	clientID, err := checkVisitRepository.FindClientID([16]byte{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func Test_VisitRepository_FindClientID_New(t *testing.T) {
	nested := new(NestedVisitRepository)
	checkVisitRepository := NewVisitRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.New().Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil)

	clientID, err := checkVisitRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func Test_VisitRepository_FindClientID_Cache(t *testing.T) {
	nested := new(NestedVisitRepository)
	checkVisitRepository := NewVisitRepository(nested, 10)
	expected := "12345"
	sessionID := uuid.New().Generate()

	nested.On("FindClientID", sessionID).Return(expected, nil).Once()
	checkVisitRepository.FindClientID(sessionID)

	clientID, err := checkVisitRepository.FindClientID(sessionID)

	assert.Empty(t, err)
	assert.Equal(t, expected, clientID)
	nested.AssertExpectations(t)
}

func Test_VisitRepository_Verify_EmptySessionID(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify_EmptyClientID(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

type NestedVisitRepository struct {
	mock.Mock
}

func (repository *NestedVisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *NestedVisitRepository) FindSessionID(clientID string) (sessionID [16]byte, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.([16]byte)

	return sessionID, args.Error(1)
}

func (repository *NestedVisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *NestedVisitRepository) Insert(visit *entities.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}
