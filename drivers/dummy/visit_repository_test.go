package dummy

import (
	"testing"

	"github.com/index0h/go-tracker/common"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/drivers/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event components.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_VisitRepository_FindClientID_Empty(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	clientID, err := checkVisitRepository.FindClientID(common.UUID{})

	assert.Empty(t, clientID)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify_EmptySessionID(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	ok, err := checkVisitRepository.Verify(common.UUID{}, "12345")

	assert.False(t, ok)
	assert.NotNil(t, err)
}

func Test_VisitRepository_Verify_EmptyClientID(t *testing.T) {
	checkVisitRepository := VisitRepository{}

	ok, err := checkVisitRepository.Verify(uuid.New().Generate(), "")

	assert.False(t, ok)
	assert.NotNil(t, err)
}
