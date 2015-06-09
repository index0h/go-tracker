package index

import (
	"testing"

	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func TestMapIndexRefreshTwo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	pushEvents := []*entities.Event{eventA, eventB}
	events := map[[16]byte]*entities.Event{eventA.EventID(): eventA, eventB.EventID(): eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(pushEvents)

	assert.Equal(t, events, testIndex.events)
}

func TestMapIndexRefreshRemoveEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entities.Event{})

	assert.Len(t, testIndex.events, 0)
}

func TestMapIndexFindAllEmpty(t *testing.T) {
	testIndex := NewMapIndex()

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestMapIndexFindAllWithData(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func TestMapIndexInsertEmpty(t *testing.T) {
	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func TestMapIndexInsertTwoEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func TestMapIndexInsertDuplicates(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func TestMapIndexDeleteEmpty(t *testing.T) {
	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func TestMapIndexDeleteEventByPointer(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestMapIndexDeleteEventByUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entities.NewEvent(eventA.EventID(), true, map[string]string{}, map[string]string{})

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventB)

	assert.Len(t, testIndex.FindAll(), 0)
}

func TestMapIndexUpdateEmptyFrom(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func TestMapIndexUpdateEmptyTo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func TestMapIndexUpdateEqual(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func TestMapIndexUpdateNotEqualUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func TestMapIndexUpdate(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entities.NewEvent(eventA.EventID(), false, map[string]string{}, map[string]string{})

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventA, testIndex.events[eventA.EventID()])

	assert.Nil(t, testIndex.Update(eventA, eventB))

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventB, testIndex.events[eventB.EventID()])
}
