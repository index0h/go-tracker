package visit

import "github.com/index0h/go-tracker/visit/entity"

type RepositoryInterface interface {
	FindByID(visitID [16]byte) (*entity.Visit, error)

	FindAll(limit int64, offset int64) ([]*entity.Visit, error)

	FindAllBySessionID(sessionID [16]byte, limit int64, offset int64) ([]*entity.Visit, error)

	FindAllByClientID(clientID string, limit int64, offset int64) ([]*entity.Visit, error)

	Insert(*entity.Visit) error

	Verify(sessionID [16]byte, clientID string) (ok bool, err error)
}

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type UUIDProviderInterface interface {
	// Generate new not empty UUID
	Generate() [16]byte

	// Converts UUID from bytes to string
	ToString([16]byte) string

	// Converts UUID from string to bytes
	ToBytes(string) [16]byte
}

