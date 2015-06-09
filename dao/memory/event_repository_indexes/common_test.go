package event_repository_indexes

import (
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func commonEventSlicesEqual(t *testing.T, first, second []*entities.Event) {
	assert.Equal(t, len(first), len(second))

	for _, eventFirst := range first {
		found := false
		for _, eventSecond := range second {
			if eventFirst.EventID() == eventSecond.EventID() {
				found = true

				assert.Equal(t, eventFirst, eventSecond)
			}
		}

		if !found {

			t.Errorf("Events slices non equal")
		}
	}
}

func commonGenerateNotFilteredEvent() *entities.Event {
	eventA, _ := entities.NewEvent(uuid.New().Generate(), true, map[string]string{}, map[string]string{})

	return eventA
}
