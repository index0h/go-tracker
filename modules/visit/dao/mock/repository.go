package mock

import (
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (repository *Repository) FindByID(visitID types.UUID) (*entity.Visit, error) {
	args := repository.Called(visitID)

	raw := args.Get(0)
	result, _ := raw.(*entity.Visit)

	return result, args.Error(1)
}

func (repository *Repository) FindAll(limit int64, offset int64) ([]*entity.Visit, error) {
	args := repository.Called(limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entity.Visit)

	return result, args.Error(1)
}

func (repository *Repository) FindAllBySessionID(
	sessionID types.UUID,
	limit int64,
	offset int64,
) ([]*entity.Visit, error) {
	args := repository.Called(sessionID, limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entity.Visit)

	return result, args.Error(1)
}

func (repository *Repository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) ([]*entity.Visit, error) {
	args := repository.Called(clientID, limit, offset)

	raw := args.Get(0)
	result, _ := raw.([]*entity.Visit)

	return result, args.Error(1)
}

func (repository *Repository) Insert(visit *entity.Visit) (err error) {
	args := repository.Called(visit)

	return args.Error(0)
}

func (repository *Repository) Verify(sessionID types.UUID, clientID string) (ok bool, err error) {
	args := repository.Called(sessionID, clientID)

	return args.Bool(0), args.Error(1)
}
