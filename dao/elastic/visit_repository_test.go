package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/entities"
	"github.com/index0h/go-tracker/dao"
	driver "github.com/olivere/elastic"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
	"fmt"
	"log"
	"os"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_Visit_Repository_FindClientID(t *testing.T) {
	client, _ := driver.NewClient(driver.SetTraceLog(log.New(os.Stdout, "logger: ", log.Lshortfile)))
	repository := NewVisitRepository(client, uuid.New())
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)



	foundClientID, err := repository.FindClientID(uuid.New().ToBytes("242f4fd4-e446-4b1a-aa25-e89ddd50d1e5"))
fmt.Printf("%+v\n\n", foundClientID)
	assert.Nil(t, err)
	assert.Equal(t, clientID, foundClientID)
}


func Test_Visit_Repository_Insert(t *testing.T) {
	client, _ := driver.NewClient()
	repository := NewVisitRepository(client, uuid.New())
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	data := map[string]string{"data": "here"}
	warnings := []string{"i'm warning"}
	timestamp := int64(15)

	indexName := repository.indexName()
	typeName := repository.typeName

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, timestamp, sessionID, clientID, data, warnings)

	err := repository.Insert(visit)

	assert.Nil(t, err)

	indexExists, err := client.IndexExists(indexName).Do()
	assert.True(t, indexExists)
	assert.Nil(t, err)

	foundRawVisit, err := client.Get().
		Index(indexName).
		Type(typeName).
		Id(uuid.New().ToString(visitID)).
		Do()

	assert.Nil(t, err)
	assert.NotEmpty(t, *foundRawVisit.Source)

	foundVisit, err := repository.byteToVisit(*foundRawVisit.Source)

	assert.Nil(t, err)
	assert.Equal(t, visit.VisitID(), foundVisit.VisitID())
	assert.Equal(t, visit.Timestamp(), foundVisit.Timestamp())
	assert.Equal(t, visit.SessionID(), foundVisit.SessionID())
	assert.Equal(t, visit.ClientID(), foundVisit.ClientID())
	assert.Equal(t, visit.Data(), foundVisit.Data())
	assert.Equal(t, visit.Warnings(), foundVisit.Warnings())
}

func Test_Visit_Repository_InsertNil(t *testing.T) {
	client, _ := driver.NewClient()
	repository := NewVisitRepository(client, uuid.New())

	err := repository.Insert(nil)

	assert.NotNil(t, err)
}
