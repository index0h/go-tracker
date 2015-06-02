package memory

import (
	"testing"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	"github.com/stretchr/testify/assert"
)

func commonEventSlicesEqual(t *testing.T, first, second []*eventEntities.Event) {
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

func commonGenerateNotFilteredEvent() *eventEntities.Event {
	uuidMaker := uuidDriver.UUID{}

	eventA, _ := eventEntities.NewEvent(uuidMaker.Generate(), true, map[string]string{}, map[string]string{})

	return eventA
}
