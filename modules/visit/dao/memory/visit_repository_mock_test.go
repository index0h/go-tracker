package memory

import (
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/mock"
)

type MockVisitRepository struct {
	mock.Mock
}

func (repository *MockVisitRepository) FindByID(visitID [16]byte) (*entities.Visit, error) {
	args := repository.Called(visitID)

	raw := args.Get(0)
	result, _ := raw.(*entities.Visit)

	return result, args.Error(1)
}

func (repository *MockVisitRepository) FindAll(limit int64, offset int64) ([]*entities.Visit, error) {
	args := repository.Called(limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entities.Visit)

	return result, args.Error(1)
}

func (repository *MockVisitRepository) FindAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	args := repository.Called(sessionID, limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entities.Visit)

	return result, args.Error(1)
}

func (repository *MockVisitRepository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	args := repository.Called(clientID, limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entities.Visit)

	return result, args.Error(1)
}

func (repository *MockVisitRepository) Insert(visit *entities.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}

func (repository *MockVisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}
