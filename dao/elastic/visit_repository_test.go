package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
)

func Test_VisitRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}

func Test_VisitRepository_FindClientID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	foundClientID, err := repository.FindClientID(sessionID)
	assert.Nil(t, err)
	assert.Equal(t, clientID, foundClientID)
}

func Test_VisitRepository_FindClientID_Empty(t *testing.T) {
	_, repository := visitRepository_CreateRepository()

	foundClientID, err := repository.FindClientID([16]byte{})
	assert.NotNil(t, err)
	assert.Empty(t, foundClientID)
}

func Test_VisitRepository_FindClientID_WrongSessionID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	foundClientID, err := repository.FindClientID(uuid.New().Generate())
	assert.Nil(t, err)
	assert.Empty(t, foundClientID)
}

func Test_VisitRepository_Verify(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(sessionID, clientID)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func Test_VisitRepository_Verify_WrongClientID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(sessionID, "Some another client ID")
	assert.Nil(t, err)
	assert.False(t, ok)
}

func Test_VisitRepository_Verify_WrongSessionID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(uuid.New().Generate(), clientID)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func Test_VisitRepository_Verify_EmptyClientID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(uuid.New().Generate(), "")
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func Test_VisitRepository_Verify_EmptySessionID(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	indexName := repository.indexName()

	_, _ = client.DeleteIndex(indexName).Do()

	visit, _ := entities.NewVisit(visitID, int64(15), sessionID, clientID, map[string]string{}, []string{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify([16]byte{}, clientID)
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func Test_VisitRepository_Insert(t *testing.T) {
	client, repository := visitRepository_CreateRepository()
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

func Test_VisitRepository_Insert_Nil(t *testing.T) {
	_, repository := visitRepository_CreateRepository()

	err := repository.Insert(nil)

	assert.NotNil(t, err)
}

func visitRepository_CreateRepository() (*driver.Client, *VisitRepository) {
	//client, _ := driver.NewClient(driver.SetTraceLog(log.New(os.Stdout, "logger: ", log.Lshortfile)))
	client, _ := driver.NewClient()
	repository, _ := NewVisitRepository(client, uuid.New())
	repository.indexPrefix = "tracker-test-"
	repository.RefreshAfterInsert = true

	return client, repository
}
