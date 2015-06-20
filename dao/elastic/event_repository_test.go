package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"time"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&EventRepository{})
}

func Test_EventRepository_NewEventRepository_EmptyClient(t *testing.T) {
	repository, err := NewEventRepository(nil, uuid.New())

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func Test_EventRepository_NewEventRepository_EmptyUUIDProvider(t *testing.T) {
	client, _ := driver.NewClient()
	repository, err := NewEventRepository(client, nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func Test_EventRepository_FindAll_Empty(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	foundEvents, err := repository.FindAll()

	assert.Nil(t, err)
	assert.Len(t, foundEvents, 0)
}

func Test_EventRepository_FindAllByVisit_Empty(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	foundEvents, err := repository.FindAllByVisit(nil)

	assert.NotNil(t, err)
	assert.Empty(t, foundEvents)
}

func Test_EventRepository_FindAllByVisit_RealVisit(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	eventA := eventRepository_GenerateEvent(map[string]string{"A": "A"})
	eventB := eventRepository_GenerateEvent(map[string]string{"B": "B"})
	eventC := eventRepository_GenerateEvent(map[string]string{"A": "B"})
	eventD := eventRepository_GenerateEvent(map[string]string{})
	eventE, _ := entities.NewEvent(uuid.New().Generate(), false, map[string]string{}, map[string]string{"Z": "Z"})

	visit := eventRepository_GenerateVisit(map[string]string{"A": "A", "B": "B"})

	events := []*entities.Event{eventA, eventB, eventD}

	repository.Insert(eventA)
	repository.Insert(eventB)
	repository.Insert(eventC)
	repository.Insert(eventD)
	repository.Insert(eventE)

	foundEvents, err := repository.FindAllByVisit(visit)

	assert.Nil(t, err)
	eventRepository_EventSlicesEqual(t, events, foundEvents)
}

func Test_EventRepository_FindAll_WithData(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	eventA := eventRepository_GenerateEvent(map[string]string{})
	eventB := eventRepository_GenerateEvent(map[string]string{})
	events := []*entities.Event{eventA, eventB}

	repository.Insert(eventA)
	repository.Insert(eventB)

	foundEvents, err := repository.FindAll()

	assert.Nil(t, err)
	eventRepository_EventSlicesEqual(t, events, foundEvents)
}

func Test_EventRepository_FindByID_Empty(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	foundEvent, err := repository.FindByID([16]byte{})

	assert.Nil(t, err)
	assert.Nil(t, foundEvent)
}

func Test_EventRepository_FindByID_NotFound(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	foundEvent, err := repository.FindByID(uuid.New().Generate())

	assert.Nil(t, err)
	assert.Nil(t, foundEvent)
}

func Test_EventRepository_FindByID_RealID(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	eventA := eventRepository_GenerateEvent(map[string]string{})
	events := []*entities.Event{eventA}

	repository.Insert(eventA)

	foundEvent, err := repository.FindByID(eventA.EventID())

	assert.Nil(t, err)
	eventRepository_EventSlicesEqual(t, events, []*entities.Event{foundEvent})
}

func Test_EventRepository_Insert_Empty(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	err := repository.Insert(nil)

	assert.NotNil(t, err)
}

func Test_EventRepository_Insert_Real(t *testing.T) {
	client, repository := eventRepository_CreateRepository()

	eventA := eventRepository_GenerateEvent(map[string]string{})

	err := repository.Insert(eventA)

	assert.Nil(t, err)

	foundRawEvent, err := client.Get().
		Index(repository.indexName).
		Type(repository.typeName).
		Id(uuid.New().ToString(eventA.EventID())).
		Refresh(true).
		Do()

	assert.Nil(t, err)

	foundEvent, err := repository.byteToEvent(*foundRawEvent.Source)

	assert.Nil(t, err)

	assert.Equal(t, eventA, foundEvent)
}

func Test_EventRepository_Update_Nil(t *testing.T) {
	_, repository := eventRepository_CreateRepository()

	err := repository.Update(nil)

	assert.NotNil(t, err)
}

func Test_EventRepository_Update_Real(t *testing.T) {
	client, repository := eventRepository_CreateRepository()

	eventA := eventRepository_GenerateEvent(map[string]string{"A": "A"})
	eventB, _ := entities.NewEvent(eventA.EventID(), false, eventA.Data(), map[string]string{"B": "B"})

	_ = repository.Insert(eventA)

	err := repository.Update(eventB)

	assert.Nil(t, err)

	foundRawEvent, err := client.Get().
		Index(repository.indexName).
		Type(repository.typeName).
		Id(uuid.New().ToString(eventB.EventID())).
		Refresh(true).
		Do()

	assert.Nil(t, err)

	foundEvent, err := repository.byteToEvent(*foundRawEvent.Source)

	assert.Nil(t, err)

	assert.Equal(t, eventB, foundEvent)
}

func eventRepository_CreateRepository() (*driver.Client, *EventRepository) {
	client, _ := driver.NewClient()
	repository, _ := NewEventRepository(client, uuid.New())
	repository.indexName = "tracker-test"

	_, _ = client.DeleteIndex(repository.indexName).Do()
	_, _ = client.CreateIndex(repository.indexName).Do()

	time.Sleep(100 * time.Millisecond)

	return client, repository
}

func eventRepository_EventSlicesEqual(t *testing.T, first, second []*entities.Event) {
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

func eventRepository_GenerateEvent(filters map[string]string) *entities.Event {
	eventA, _ := entities.NewEvent(uuid.New().Generate(), true, map[string]string{"data": "here"}, filters)

	return eventA
}

func eventRepository_GenerateVisit(data map[string]string) *entities.Visit {
	visit, _ := entities.NewVisit(uuid.New().Generate(), int64(0), uuid.New().Generate(), "", data, []string{})

	return visit
}
