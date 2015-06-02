package memory

import (
	"testing"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	"github.com/stretchr/testify/assert"
)

func TestListIndexRefreshTwo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func TestListIndexRefreshRemoveEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	newEvents := []*eventEntities.Event{}
	testIndex.Refresh(newEvents)

	assert.Len(t, testIndex.events, 0)
}

func TestListIndexRefreshCopies(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB, eventA, eventB}
	eventsClean := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, eventsClean, testIndex.events)
}

func TestListIndexFindAllEmpty(t *testing.T) {
	testIndex := NewListIndex()

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestListIndexFindAllWithData(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func TestListIndexInsertEmpty(t *testing.T) {
	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func TestListIndexInsertTwoEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func TestListIndexInsertDuplicates(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func TestListIndexDeleteEmpty(t *testing.T) {
	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func TestListIndexDeleteEventByPointer(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestListIndexDeleteEventByUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := eventEntities.NewEvent(eventA.EventID(), true, map[string]string{}, map[string]string{})

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventB)

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestListIndexUpdateEmptyFrom(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func TestListIndexUpdateEmptyTo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func TestListIndexUpdateEqual(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func TestListIndexUpdateNotEqualUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func TestListIndexUpdate(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := eventEntities.NewEvent(eventA.EventID(), false, map[string]string{}, map[string]string{})

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventA, testIndex.events[0])

	assert.Nil(t, testIndex.Update(eventA, eventB))

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventB, testIndex.events[0])

}
