package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_VisitRepository_FindClientID(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	clientID, err := checkVisitRepository.FindClientID(uuid.New().Generate())

	assert.Empty(t, clientID)
	assert.Nil(t, err)
}

func Test_VisitRepository_FindClientID_Empty(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	clientID, err := checkVisitRepository.FindClientID([16]byte{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "12345")

	assert.True(t, ok)
	assert.Nil(t, err)
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

func Test_VisitRepository_Insert(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	data := map[string]string{"data": "here"}
	warnings := []string{"i'm warning"}
	timestamp := int64(15)

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, data, warnings)

	err := checkVisitRepository.Insert(visit)

	assert.Nil(t, err)
}

func Test_VisitRepository_Insert_Nil(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	err := checkVisitRepository.Insert(nil)

	assert.NotNil(t, err)
}
