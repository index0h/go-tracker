package memory

import (
	"errors"

	"github.com/golang/groupcache/lru"
	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct {
	nested visit.RepositoryInterface
	cache  *lru.Cache
}

func NewRepository(nested visit.RepositoryInterface, maxEntries int) (*Repository, error) {
	if nested == nil {
		return nil, errors.New("Empty nested is not allowed")
	}

	return &Repository{
		nested: nested,
		cache:  lru.New(maxEntries),
	}, nil
}

func (repository *Repository) FindByID(visitID types.UUID) (result *entity.Visit, err error) {
	if visitID.IsEmpty() {
		return result, errors.New("Empty visitID is not allowed")
	}

	return repository.nested.FindByID(visitID)
}

func (repository *Repository) FindAll(limit int64, offset int64) (result []*entity.Visit, err error) {
	return repository.nested.FindAll(limit, offset)
}

func (repository *Repository) FindAllBySessionID(
	sessionID types.UUID,
	limit int64,
	offset int64,
) (result []*entity.Visit, err error) {
	if sessionID.IsEmpty() {
		return result, errors.New("Empty sessionID is not allowed")
	}

	return repository.nested.FindAllBySessionID(sessionID, limit, offset)
}

func (repository *Repository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) (result []*entity.Visit, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	return repository.nested.FindAllByClientID(clientID, limit, offset)
}

// Save visit to cache and run nested save
func (repository *Repository) Insert(visit *entity.Visit) (err error) {
	if visit == nil {
		return errors.New("visit must be not nil")
	}

	if visit.ClientID() != "" {
		repository.cache.Add(visit.SessionID(), visit.ClientID())
	}

	return repository.nested.Insert(visit)
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
// If sessionID or clientID not found it'll run nested repository and cache result (if its ok)
func (repository *Repository) Verify(sessionID types.UUID, clientID string) (ok bool, err error) {
	if sessionID.IsEmpty() {
		return false, errors.New("Empty sessionID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	if foundRaw, ok := repository.cache.Get(sessionID); ok {
		foundClientID, _ := foundRaw.(string)

		if foundClientID == clientID {
			return true, nil
		} else {
			return false, errors.New("sessionID registered by another clientID")
		}
	}

	if ok, err := repository.nested.Verify(sessionID, clientID); !ok || (err != nil) {
		return ok, err
	}

	repository.cache.Add(clientID, sessionID)

	return true, nil
}
