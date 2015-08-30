package memory

import (
	"errors"

	"github.com/golang/groupcache/lru"
	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
)

type Repository struct {
	nested          visit.RepositoryInterface
	sessionToClient *lru.Cache
	clientToSession *lru.Cache
}

func NewRepository(nested visit.RepositoryInterface, maxEntries int) (*Repository, error) {
	if nested == nil {
		return nil, errors.New("Empty nested is not allowed")
	}

	return &Repository{
		nested:          nested,
		sessionToClient: lru.New(int(maxEntries / 2)),
		clientToSession: lru.New(int(maxEntries / 2)),
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

	repository.sessionToClient.Add(visit.SessionID(), visit.ClientID())

	if visit.ClientID() != "" {
		repository.clientToSession.Add(visit.ClientID(), visit.SessionID())
	}

	return repository.nested.Insert(visit)
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
// If sessionID or clientID not found it'll run nested repository and cache result (if its ok)
func (repository *Repository) Verify(sessionID types.UUID, clientID string) (ok bool, err error) {
	if sessionID.IsEmpty() {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	var (
		foundRaw       interface{}
		foundSessionID types.UUID
		foundClientID  string
	)

	if foundRaw, ok = repository.sessionToClient.Get(sessionID); ok {
		foundClientID, _ = foundRaw.(string)

		if foundClientID == clientID {
			repository.clientToSession.Add(clientID, sessionID)

			return true, nil
		}

		if foundClientID != "" {
			return false, errors.New("sessionID registered by another clientID")
		}
	}

	if foundRaw, ok = repository.clientToSession.Get(clientID); ok {
		foundSessionID, _ = foundRaw.(types.UUID)

		if foundSessionID == sessionID {
			repository.sessionToClient.Add(sessionID, clientID)

			return true, nil
		}

		if !foundSessionID.IsEmpty() {
			return false, errors.New("sessionID registered by another clientID")
		}
	}

	ok, err = repository.nested.Verify(sessionID, clientID)
	if !ok || err != nil {
		return ok, err
	}

	repository.sessionToClient.Add(sessionID, clientID)
	repository.clientToSession.Add(clientID, sessionID)

	return true, nil
}
