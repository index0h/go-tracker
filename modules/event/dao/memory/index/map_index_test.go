package index

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_MapIndex_Refresh_Two(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	pushEvents := []*entity.Event{eventA, eventB}
	events := map[[16]byte]*entity.Event{eventA.EventID(): eventA, eventB.EventID(): eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(pushEvents)

	assert.Equal(t, events, testIndex.events)
}

func Test_MapIndex_Refresh_RemoveEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entity.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entity.Event{})

	assert.Len(t, testIndex.events, 0)
}

func Test_MapIndex_FindAll_Empty(t *testing.T) {
	testIndex := NewMapIndex()

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_MapIndex_FindAll_WithData(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entity.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(events)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func Test_EventRepository_FindByID_Empty(t *testing.T) {
	testIndex := NewMapIndex()

	foundEvent, err := testIndex.FindByID([16]byte{})

	assert.NotNil(t, err)
	assert.Nil(t, foundEvent)
}

func Test_EventRepository_FindByID_NotFound(t *testing.T) {
	testIndex := NewMapIndex()

	foundEvent, err := testIndex.FindByID(uuid.New().Generate())

	assert.Nil(t, err)
	assert.Nil(t, foundEvent)
}

func Test_EventRepository_FindByID_RealID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entity.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Refresh(events)

	foundEvent, err := testIndex.FindByID(eventA.EventID())

	assert.Nil(t, err)
	assert.Equal(t, eventA.EventID(), foundEvent.EventID())
}

func Test_MapIndex_Insert_Nil(t *testing.T) {
	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func Test_MapIndex_Insert_TwoEvents(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entity.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func Test_MapIndex_Insert_Duplicates(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	events := []*entity.Event{eventA, eventB}

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventB)

	commonEventSlicesEqual(t, events, testIndex.FindAll())
}

func Test_MapIndex_Delete_Nil(t *testing.T) {
	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func Test_MapIndex_Delete_EventByPointer(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_MapIndex_Delete_EventByUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entity.NewEvent(eventA.EventID(), true, types.Hash{}, types.Hash{})

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventB)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_MapIndex_Delete_NotPresent(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.FindAll(), 0)
}

func Test_MapIndex_Update_EmptyFrom(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func Test_MapIndex_Update_EmptyTo(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func Test_MapIndex_Update_Equal(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func Test_MapIndex_Update_NotEqualUUID(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB := commonGenerateNotFilteredEvent()

	testIndex := NewMapIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func Test_MapIndex_Update(t *testing.T) {
	eventA := commonGenerateNotFilteredEvent()
	eventB, _ := entity.NewEvent(eventA.EventID(), false, types.Hash{}, types.Hash{})

	testIndex := NewMapIndex()
	testIndex.Insert(eventA)

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventA, testIndex.events[eventA.EventID()])

	assert.Nil(t, testIndex.Update(eventA, eventB))

	assert.Len(t, testIndex.events, 1)
	assert.Equal(t, eventB, testIndex.events[eventB.EventID()])
}
