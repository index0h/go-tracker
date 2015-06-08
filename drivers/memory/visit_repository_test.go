package memory

import (
	"testing"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/drivers/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInterface(t *testing.T) {
	func(event components.VisitRepositoryInterface) {}(&VisitRepository{})
}

func TestFindClientIDEmpty(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	clientID, err := checkVisitRepository.FindClientID(common.UUID{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func TestFindClientIDNew(t *testing.T) {
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

func TestFindClientIDCache(t *testing.T) {
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

func TestVerifyEmptySessionID(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify(common.UUID{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestVerifyEmptyClientID(t *testing.T) {
	checkVisitRepository := NewVisitRepository(new(NestedVisitRepository), 10)

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

type NestedVisitRepository struct {
	mock.Mock
}

func (repository *NestedVisitRepository) FindClientID(sessionID common.UUID) (clientID string, err error) {
	args := repository.Called(sessionID)

	return args.String(0), args.Error(1)
}

func (repository *NestedVisitRepository) FindSessionID(clientID string) (sessionID common.UUID, err error) {
	args := repository.Called(clientID)

	raw := args.Get(0)
	sessionID, _ = raw.(common.UUID)

	return sessionID, args.Error(1)
}

func (repository *NestedVisitRepository) Verify(sessionID common.UUID, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}

func (repository *NestedVisitRepository) Insert(visit *entities.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}
