package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_VisitRepository_FindClientID_Empty(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	clientID, err := checkVisitRepository.FindClientID([16]byte{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify_EmptySessionID(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	ok, err := checkVisitRepository.Verify([16]byte{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify_EmptyClientID(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}
