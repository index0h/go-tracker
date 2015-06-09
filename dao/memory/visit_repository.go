package memory

import (
	"errors"

	"github.com/golang/groupcache/lru"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type VisitRepository struct {
	nested          dao.VisitRepositoryInterface
	sessionToClient *lru.Cache
	clientToSession *lru.Cache
}

func NewVisitRepository(nested dao.VisitRepositoryInterface, maxEntries int) *VisitRepository {
	return &VisitRepository{
		nested:          nested,
		sessionToClient: lru.New(int(maxEntries / 2)),
		clientToSession: lru.New(int(maxEntries / 2)),
	}
}

// Find clientID by sessionID. If it's not present in cache - will try to find by nested repository and cache result
func (repository *VisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	if sessionID == [16]byte{} {
		return clientID, errors.New("Empty sessioID is not allowed")
	}

	if rawFound, ok := repository.sessionToClient.Get(sessionID); ok {
		clientID, _ = rawFound.(string)

		return clientID, nil
	}

	clientID, err = repository.nested.FindClientID(sessionID)

	if err == nil {
		repository.sessionToClient.Add(sessionID, clientID)

		if clientID != "" {
			repository.clientToSession.Add(clientID, sessionID)
		}
	}

	return clientID, err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
// If sessionID or clientID not found it'll run nested repository and cache result (if its ok)
func (repository *VisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	if sessionID == [16]byte{} {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	var (
		foundRaw       interface{}
		foundSessionID [16]byte
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
		foundSessionID, _ = foundRaw.([16]byte)

		if foundSessionID == sessionID {
			repository.sessionToClient.Add(sessionID, clientID)

			return true, nil
		}

		if foundSessionID != [16]byte{} {
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

// Save visit to cache and run nested save
func (repository *VisitRepository) Insert(visit *entities.Visit) (err error) {
	if visit == nil {
		return errors.New("visit must be not nil")
	}

	repository.sessionToClient.Add(visit.SessionID(), visit.ClientID())

	if visit.ClientID() != "" {
		repository.clientToSession.Add(visit.ClientID(), visit.SessionID())
	}

	return repository.nested.Insert(visit)
}
