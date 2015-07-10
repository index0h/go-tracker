package index

import (
	"testing"

	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func Test_ListIndex_Refresh_Two(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func Test_ListIndex_Refresh_RemoveEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	newEvents := []*entities.Event{}
	testIndex.Refresh(newEvents)

	assert.Len(t, testIndex.events, 0)
}

func Test_ListIndex_Refresh_Copies(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB, eventA, eventB}
	eventsClean := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, eventsClean, testIndex.events)
}

func Test_ListIndex_FindAll_Empty(t *testing.T) {
	testIndex := NewListIndex()

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_ListIndex_FindAll_WithData(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func Test_ListIndex_InsertEmpty(t *testing.T) {
	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func Test_ListIndex_Insert_TwoEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func Test_ListIndex_Insert_Duplicates(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entities.Event{eventA, eventB}

	testIndex := NewListIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.events)
}

func Test_ListIndex_Delete_Empty(t *testing.T) {
	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func Test_ListIndex_Delete_EventByPointer(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_ListIndex_Delete_EventByUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entities.NewEvent(eventA.EventID(), true, entities.Hash{}, entities.Hash{})

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventB)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_ListIndex_Delete_NotFound(t *testing.T) {
	event := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	testIndex.Delete(event)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_ListIndex_Update_EmptyFrom(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func Test_ListIndex_Update_EmptyTo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func Test_ListIndex_Update_Equal(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func Test_ListIndex_Update_NotEqualUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	testIndex := NewListIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func Test_ListIndex_Update(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entities.NewEvent(eventA.EventID(), false, entities.Hash{}, entities.Hash{})

	testIndex := NewListIndex()
	testIndex.Insert(eventA)

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventA, testIndex.events[0])

	assert.Nil(t, testIndex.Update(eventA, eventB))

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventB, testIndex.events[0])

}
