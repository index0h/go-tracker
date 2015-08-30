package index

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func commonEventSlicesEqual(t *testing.T, first, second []*entity.Event) {
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

func commonGenerateNotFilteredEvent() *entity.Event {
	eventA, _ := entity.NewEvent(uuid.New().Generate(), true, types.Hash{}, types.Hash{})

	return eventA
}
