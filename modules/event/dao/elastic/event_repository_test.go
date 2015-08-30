package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/modules/event"
	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	driver "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestRepository_Interface(t *testing.T) {
	func(event event.RepositoryInterface) {}(&Repository{})
}

func TestRepository_NewRepository_EmptyClient(t *testing.T) {
	repository, err := NewRepository(nil, uuid.New())

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestRepository_NewRepository_EmptyUUIDProvider(t *testing.T) {
	client, _ := driver.NewClient()
	repository, err := NewRepository(client, nil)

	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestRepository_FindAll(t *testing.T) {
	_, repository := Repository_CreateRepository()
	foundEvents, err := repository.FindAll(0, 0)

	assert.Nil(t, err)
	assert.Len(t, foundEvents, 0)
}

func TestRepository_FindAllByVisit_Empty(t *testing.T) {
	_, repository := Repository_CreateRepository()

	foundEvents, err := repository.FindAllByFields(nil)

	assert.NotNil(t, err)
	assert.Empty(t, foundEvents)
}

func TestRepository_FindAllByFields_RealVisit(t *testing.T) {
	_, repository := Repository_CreateRepository()

	eventA := Repository_GenerateEvent(types.Hash{"A": "A"})
	eventB := Repository_GenerateEvent(types.Hash{"B": "B"})
	eventC := Repository_GenerateEvent(types.Hash{"A": "B"})
	eventD := Repository_GenerateEvent(types.Hash{})
	eventE, _ := entity.NewEvent(uuid.New().Generate(), false, types.Hash{}, types.Hash{"Z": "Z"})

	events := []*entity.Event{eventA, eventB, eventD}

	repository.Insert(eventA)
	repository.Insert(eventB)
	repository.Insert(eventC)
	repository.Insert(eventD)
	repository.Insert(eventE)

	foundEvents, err := repository.FindAllByFields(types.Hash{"A": "A", "B": "B"})

	assert.Nil(t, err)
	Repository_EventSlicesEqual(t, events, foundEvents)
}

func TestRepository_FindAllByFields_NoEventsForVisit(t *testing.T) {
	_, repository := Repository_CreateRepository()

	foundEvents, err := repository.FindAllByFields(types.Hash{"A": "A", "B": "B"})

	assert.Nil(t, err)
	assert.Empty(t, foundEvents)
}

func TestRepository_FindAll_WithData(t *testing.T) {
	_, repository := Repository_CreateRepository()

	eventA := Repository_GenerateEvent(types.Hash{})
	eventB := Repository_GenerateEvent(types.Hash{})
	events := []*entity.Event{eventA, eventB}

	repository.Insert(eventA)
	repository.Insert(eventB)

	foundEvents, err := repository.FindAll(0, 0)

	assert.Nil(t, err)
	Repository_EventSlicesEqual(t, events, foundEvents)
}

func TestRepository_FindByID_Empty(t *testing.T) {
	_, repository := Repository_CreateRepository()

	foundEvent, err := repository.FindByID([16]byte{})

	assert.NotNil(t, err)
	assert.Nil(t, foundEvent)
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	_, repository := Repository_CreateRepository()

	foundEvent, err := repository.FindByID(uuid.New().Generate())

	assert.Nil(t, err)
	assert.Nil(t, foundEvent)
}

func TestRepository_FindByID_RealID(t *testing.T) {
	_, repository := Repository_CreateRepository()

	eventA := Repository_GenerateEvent(types.Hash{})
	events := []*entity.Event{eventA}

	repository.Insert(eventA)

	foundEvent, err := repository.FindByID(eventA.EventID())

	assert.Nil(t, err)
	Repository_EventSlicesEqual(t, events, []*entity.Event{foundEvent})
}

func TestRepository_Insert_Empty(t *testing.T) {
	_, repository := Repository_CreateRepository()

	err := repository.Insert(nil)

	assert.NotNil(t, err)
}

func TestRepository_Insert_Real(t *testing.T) {
	client, repository := Repository_CreateRepository()

	eventA := Repository_GenerateEvent(types.Hash{})

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

func TestRepository_Update_Nil(t *testing.T) {
	_, repository := Repository_CreateRepository()

	err := repository.Update(nil)

	assert.NotNil(t, err)
}

func TestRepository_Update_Real(t *testing.T) {
	client, repository := Repository_CreateRepository()

	eventA := Repository_GenerateEvent(types.Hash{"A": "A"})
	eventB, _ := entity.NewEvent(eventA.EventID(), false, eventA.Fields(), types.Hash{"B": "B"})

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

func Repository_CreateRepository() (*driver.Client, *Repository) {
	client, _ := driver.NewClient()
	repository, _ := NewRepository(client, uuid.New())
	repository.indexName = "test-tracker"

	_, _ = client.DeleteIndex(repository.indexName).Do()
	_, _ = client.CreateIndex(repository.indexName).Do()

	time.Sleep(100 * time.Millisecond)

	return client, repository
}

func Repository_EventSlicesEqual(t *testing.T, first, second []*entity.Event) {
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

func Repository_GenerateEvent(filters types.Hash) *entity.Event {
	eventA, _ := entity.NewEvent(uuid.New().Generate(), true, types.Hash{"data": "here"}, filters)

	return eventA
}
