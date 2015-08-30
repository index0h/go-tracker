package index

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_FilterIndex_Refresh_Two(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A1", "B": "B", "C": "C1"})
	eventB := filterIndexGenerateEvent(types.Hash{"B": "B", "C": "C2", "D": "D"})

	pushEvents := []*entity.Event{eventA, eventB}

	testIndex := NewFilterIndex()
	testIndex.Refresh(pushEvents)

	assert.NotEmpty(t, testIndex.events["A"])
	assert.NotEmpty(t, testIndex.events["B"])
	assert.NotEmpty(t, testIndex.events["C"])
	assert.NotEmpty(t, testIndex.events["D"])

	assert.Len(t, testIndex.events["A"]["A1"], 1)
	assert.Len(t, testIndex.events["B"]["B"], 2)
	assert.Len(t, testIndex.events["C"]["C1"], 1)
	assert.Len(t, testIndex.events["C"]["C2"], 1)
	assert.Len(t, testIndex.events["D"]["D"], 1)

	assert.Equal(t, uint(3), testIndex.events["A"]["A1"][eventA])
	assert.Equal(t, uint(3), testIndex.events["B"]["B"][eventA])
	assert.Equal(t, uint(3), testIndex.events["C"]["C1"][eventA])

	assert.Equal(t, uint(3), testIndex.events["B"]["B"][eventB])
	assert.Equal(t, uint(3), testIndex.events["C"]["C2"][eventB])
	assert.Equal(t, uint(3), testIndex.events["D"]["D"][eventB])
}

func Test_FilterIndex_Refresh_RemoveEvents(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A1"})
	eventB := filterIndexGenerateEvent(types.Hash{"A": "A2"})

	events := []*entity.Event{eventA, eventB}

	testIndex := NewFilterIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entity.Event{})

	assert.Len(t, testIndex.events, 0)
}

func Test_FilterIndex_Refresh_Disabled(t *testing.T) {
	event, _ := entity.NewEvent(uuid.New().Generate(), false, types.Hash{}, types.Hash{"A": "A"})

	events := []*entity.Event{event}

	testIndex := NewFilterIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entity.Event{})

	assert.Len(t, testIndex.events, 0)
}

func Test_FilterIndex_FindAllByFields_Empty(t *testing.T) {
	testIndex := NewFilterIndex()

	event, err := testIndex.FindAllByFields(nil)

	assert.Nil(t, event)
	assert.NotNil(t, err)
}

func Test_FilterIndex_FindAllByFields_WithData(t *testing.T) {
	eventA1 := filterIndexGenerateEvent(types.Hash{"A": "A1"})
	eventA2 := filterIndexGenerateEvent(types.Hash{"A": "A2"})
	eventB := filterIndexGenerateEvent(types.Hash{"B": "B"})
	eventA1B := filterIndexGenerateEvent(types.Hash{"A": "A1", "B": "B"})
	eventA2B := filterIndexGenerateEvent(types.Hash{"A": "A2", "B": "B"})

	listA1 := []*entity.Event{eventA1}
	listA2 := []*entity.Event{eventA2}
	listB := []*entity.Event{eventB}
	listA1B := []*entity.Event{eventA1, eventB, eventA1B}
	listA2B := []*entity.Event{eventA2, eventB, eventA2B}

	testIndex := NewFilterIndex()

	testIndex.Insert(eventA1)
	testIndex.Insert(eventA2)
	testIndex.Insert(eventB)
	testIndex.Insert(eventA1B)
	testIndex.Insert(eventA2B)

	combinations := []filterIndexFixture{
		{fields: types.Hash{"A": "A1"}, events: listA1},
		{fields: types.Hash{"A": "A2"}, events: listA2},
		{fields: types.Hash{"B": "B"}, events: listB},
		{fields: types.Hash{"A": "A1", "B": "B"}, events: listA1B},
		{fields: types.Hash{"A": "A2", "B": "B"}, events: listA2B},
	}

	for _, fixture := range combinations {
		foundByVisit, _ := testIndex.FindAllByFields(fixture.fields)

		commonEventSlicesEqual(t, fixture.events, foundByVisit)
	}
}

func Test_FilterIndex_Insert_Empty(t *testing.T) {
	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func Test_FilterIndex_Insert_TwoEvents(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A1"})
	eventB := filterIndexGenerateEvent(types.Hash{"A": "A2"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)

	assert.NotEmpty(t, testIndex.events["A"])

	assert.Len(t, testIndex.events["A"], 2)

	assert.Len(t, testIndex.events["A"]["A1"], 1)
	assert.Len(t, testIndex.events["A"]["A2"], 1)

	assert.Equal(t, uint(1), testIndex.events["A"]["A1"][eventA])
	assert.Equal(t, uint(1), testIndex.events["A"]["A2"][eventB])
}

func Test_FilterIndex_Insert_Duplicates(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A1"})
	eventB := filterIndexGenerateEvent(types.Hash{"A": "A2"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventB)

	assert.NotEmpty(t, testIndex.events["A"])

	assert.Len(t, testIndex.events["A"], 2)

	assert.Len(t, testIndex.events["A"]["A1"], 1)
	assert.Len(t, testIndex.events["A"]["A2"], 1)

	assert.Equal(t, uint(1), testIndex.events["A"]["A1"][eventA])
	assert.Equal(t, uint(1), testIndex.events["A"]["A2"][eventB])
}

func Test_FilterIndex_Delete_Empty(t *testing.T) {
	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func Test_FilterIndex_Delete_EventByPointer(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Empty(t, testIndex.events)
}

func Test_FilterIndex_Delete_SameKeyValue(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})
	eventB := filterIndexGenerateEvent(types.Hash{"A": "A"})
	eventC := filterIndexGenerateEvent(types.Hash{"A": "B"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventC)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.events, 1)
}

func Test_FilterIndex_Update_EmptyFrom(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func Test_FilterIndex_Update_EmptyTo(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func Test_FilterIndex_Update_Equal(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func Test_FilterIndex_Update_NotEqualUUID(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})
	eventB := filterIndexGenerateEvent(types.Hash{"B": "B"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func Test_FilterIndex_Update(t *testing.T) {
	eventA := filterIndexGenerateEvent(types.Hash{"A": "A"})
	eventB, _ := entity.NewEvent(eventA.EventID(), false, types.Hash{}, types.Hash{"B": "B"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)

	assert.Empty(t, testIndex.events["B"])

	assert.NotEmpty(t, testIndex.events["A"])
	assert.Len(t, testIndex.events["A"], 1)
	assert.Len(t, testIndex.events["A"]["A"], 1)
	assert.Equal(t, uint(1), testIndex.events["A"]["A"][eventA])

	assert.Nil(t, testIndex.Update(eventA, eventB))

	assert.NotEmpty(t, testIndex.events["B"])
	assert.Len(t, testIndex.events["B"], 1)
	assert.Len(t, testIndex.events["B"]["B"], 1)
	assert.Equal(t, uint(1), testIndex.events["B"]["B"][eventB])

	assert.Empty(t, testIndex.events["A"])
}

func filterIndexGenerateEvent(filters types.Hash) *entity.Event {
	event, _ := entity.NewEvent(uuid.New().Generate(), true, types.Hash{}, filters)

	return event
}

type filterIndexFixture struct {
	events []*entity.Event
	fields types.Hash
}
