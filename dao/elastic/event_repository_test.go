package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	driver "github.com/olivere/elastic"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event dao.EventRepositoryInterface) {}(&EventRepository{})
}

func eventRepository_CreateRepository() (*driver.Client, *EventRepository) {
	//client, _ := driver.NewClient(driver.SetTraceLog(log.New(os.Stdout, "logger: ", log.Lshortfile)))
	client, _ := driver.NewClient()
	repository, _ := NewEventRepository(client, uuid.New())
	repository.indexName = "tracker-test"

	return client, repository
}
