package memory

import (
	"testing"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
	"github.com/stretchr/testify/assert"
)

func TestFilterIndexRefreshTwo(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1", "B": "B", "C": "C1"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B", "C": "C2", "D": "D"})

	pushEvents := []*eventEntities.Event{eventA, eventB}

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

func TestFilterIndexRefreshRemoveEvents(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventB := filterIndexGenerateEvent(map[string]string{"A": "A2"})

	events := []*eventEntities.Event{eventA, eventB}

	testIndex := NewFilterIndex()
	testIndex.Refresh(events)

	testIndex.Refresh([]*eventEntities.Event{})

	assert.Len(t, testIndex.events, 0)
}

func TestFilterIndexFindAllByVisitEmpty(t *testing.T) {
	testIndex := NewFilterIndex()

	event, error := testIndex.FindAllByVisit(nil)

	assert.Nil(t, event)
	assert.NotNil(t, error)
}

func TestFilterIndexFindAllByNotFoundVisit(t *testing.T) {
	testIndex := NewFilterIndex()

	visit := filterIndexGenerateVisit(map[string]string{"A": "A"})

	event, error := testIndex.FindAllByVisit(visit)

	assert.Nil(t, event)
	assert.Nil(t, error)
}

func TestFilterIndexFindAllByVisitWithData(t *testing.T) {
	eventA1 := filterIndexGenerateEvent(map[string]string{"A": "A1"})
	eventA2 := filterIndexGenerateEvent(map[string]string{"A": "A2"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B"})
	eventA1B := filterIndexGenerateEvent(map[string]string{"A": "A1", "B": "B"})
	eventA2B := filterIndexGenerateEvent(map[string]string{"A": "A2", "B": "B"})

	listA1 := []*eventEntities.Event{eventA1}
	listA2 := []*eventEntities.Event{eventA2}
	listB := []*eventEntities.Event{eventB}
	listA1B := []*eventEntities.Event{eventA1, eventB, eventA1B}
	listA2B := []*eventEntities.Event{eventA2, eventB, eventA2B}

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

func TestFilterIndexInsertEmpty(t *testing.T) {
	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Insert(nil))
}

func TestFilterIndexInsertTwoEvents(t *testing.T) {
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

func TestFilterIndexInsertDuplicates(t *testing.T) {
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

func TestFilterIndexDeleteEmpty(t *testing.T) {
	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Delete(nil))
}

func TestFilterIndexDeleteEventByPointer(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()
	testIndex.Insert(eventA)

	testIndex.Delete(eventA)

	assert.Empty(t, testIndex.events)
}

func TestFilterIndexUpdateEmptyFrom(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(nil, eventA))
}

func TestFilterIndexUpdateEmptyTo(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, nil))
}

func TestFilterIndexUpdateEqual(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventA))
}

func TestFilterIndexUpdateNotEqualUUID(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventB := filterIndexGenerateEvent(map[string]string{"B": "B"})

	testIndex := NewFilterIndex()

	assert.NotNil(t, testIndex.Update(eventA, eventB))
}

func TestFilterIndexUpdate(t *testing.T) {
	eventA := filterIndexGenerateEvent(map[string]string{"A": "A"})
	eventB, _ := eventEntities.NewEvent(eventA.EventID(), false, map[string]string{}, map[string]string{"B": "B"})

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

func filterIndexGenerateEvent(filters map[string]string) *eventEntities.Event {
	uuidMaker := uuidDriver.UUID{}

	eventA, _ := eventEntities.NewEvent(uuidMaker.Generate(), true, map[string]string{}, filters)

	return eventA
}

func filterIndexGenerateVisit(data map[string]string) *visitEntities.Visit {
	uuidMaker := uuidDriver.UUID{}

	visit, _ := visitEntities.NewVisit(uuidMaker.Generate(), int64(0), uuidMaker.Generate(), "", data, []string{})
	return visit
}

type filterIndexFixture struct {
	events []*eventEntities.Event
	fields map[string]string
}
