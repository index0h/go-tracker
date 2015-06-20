package index

import (
	"testing"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	"github.com/stretchr/testify/assert"
)

func Test_FilterIndex_Refresh_Two(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1", "B": "B", "C": "C1"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B", "C": "C2", "D": "D"})

	pushEvents := []*entities.Event{eventA, eventB}

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
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventB := filterIndexGenerateEvent(map[string]string{"A": "A2"})

	events := []*entities.Event{eventA, eventB}

	testIndex := NewFilterIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entities.Event{})

	assert.Len(t, testIndex.events, 0)
}

func Test_FilterIndex_Refresh_Disabled(t *testing.T) {
	event, _ := entities.NewEvent(uuid.New().Generate(), false, map[string]string{}, map[string]string{"A": "A"})

	events := []*entities.Event{event}

	testIndex := NewFilterIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*entities.Event{})

	assert.Len(t, testIndex.events, 0)
}

func Test_FilterIndex_FindAllByVisit_Empty(t *testing.T) {
	testIndex := NewFilterIndex()

	event, err := testIndex.FindAllByVisit(nil)

	assert.Nil(t, event)
	assert.NotNil(t, err)
}

func TestFilterIndexFindAllByNotFoundVisit(t *testing.T) {
	testIndex := NewFilterIndex()

	visit := filterIndexGenerateVisit(map[string]string{"A": "A"})

	event, err := testIndex.FindAllByVisit(visit)

	assert.Nil(t, event)
	assert.Nil(t, err)
}

func Test_FilterIndex_FindAllByVisit_WithData(t *testing.T) {
	eventA1 := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventA2 := filterIndexGenerateEvent(map[string]string{"A": "A2"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B"})
	eventA1B := filterIndexGenerateEvent(map[string]string{"A": "A1", "B": "B"})
	eventA2B := filterIndexGenerateEvent(map[string]string{"A": "A2", "B": "B"})

	listA1 := []*entities.Event{eventA1}
	listA2 := []*entities.Event{eventA2}
	listB := []*entities.Event{eventB}
	listA1B := []*entities.Event{eventA1, eventB, eventA1B}
	listA2B := []*entities.Event{eventA2, eventB, eventA2B}

	testIndex := NewFilterIndex()

	testIndex.Insert(eventA1)
	testIndex.Insert(eventA2)
	testIndex.Insert(eventB)
	testIndex.Insert(eventA1B)
	testIndex.Insert(eventA2B)

	combinations := []filterIndexFixture{
		filterIndexFixture{fields: map[string]string{"A": "A1"}, events: listA1},
		filterIndexFixture{fields: map[string]string{"A": "A2"}, events: listA2},
		filterIndexFixture{fields: map[string]string{"B": "B"}, events: listB},
		filterIndexFixture{fields: map[string]string{"A": "A1", "B": "B"}, events: listA1B},
		filterIndexFixture{fields: map[string]string{"A": "A2", "B": "B"}, events: listA2B},
	}

	for _, fixture := range combinations {
		visit := filterIndexGenerateVisit(fixture.fields)

		foundByVisit, _ := testIndex.FindAllByVisit(visit)

		commonEventSlicesEqual(t, fixture.events, foundByVisit)
	}
}

func Test_FilterIndex_Insert_Empty(t *testing.T) {
	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func Test_FilterIndex_Insert_TwoEvents(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventB := filterIndexGenerateEvent(map[string]string{"A": "A2"})

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
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventB := filterIndexGenerateEvent(map[string]string{"A": "A2"})

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
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Empty(t, testIndex.events)
}

func Test_FilterIndex_Delete_SameKeyValue(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventB := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventC := filterIndexGenerateEvent(map[string]string{"A": "B"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)
	testIndex.Insert(eventB)
	testIndex.Insert(eventC)

	testIndex.Delete(eventA)

	assert.Len(t, testIndex.events, 1)
}

func Test_FilterIndex_Update_EmptyFrom(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func Test_FilterIndex_Update_EmptyTo(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func Test_FilterIndex_Update_Equal(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func Test_FilterIndex_Update_NotEqualUUID(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func Test_FilterIndex_Update(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventB, _ := entities.NewEvent(eventA.EventID(), false, map[string]string{}, map[string]string{"B": "B"})

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

func filterIndexGenerateEvent(filters map[string]string) *entities.Event {
	event, _ := entities.NewEvent(uuid.New().Generate(), true, map[string]string{}, filters)

	return event
}

func filterIndexGenerateVisit(data map[string]string) *entities.Visit {
	visit, _ := entities.NewVisit(uuid.New().Generate(), int64(0), uuid.New().Generate(), "", data, []string{})
	return visit
}

type filterIndexFixture struct {
	events []*entities.Event
	fields map[string]string
}
