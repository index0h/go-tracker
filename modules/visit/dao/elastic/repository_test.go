package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/types"
	"github.com/index0h/go-tracker/share/uuid"
	driver "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Interface(t *testing.T) {
	func(event visit.RepositoryInterface) {}(&Repository{})
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

func TestRepository_Verify(t *testing.T) {
	_, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	visit, _ := entity.NewVisit(visitID, int64(15), sessionID, clientID, types.Hash{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(sessionID, clientID)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestRepository_Verify_WrongClientID(t *testing.T) {
	_, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	visit, _ := entity.NewVisit(visitID, int64(15), sessionID, clientID, types.Hash{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(sessionID, "Some another client ID")
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRepository_Verify_WrongSessionID(t *testing.T) {
	_, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	visit, _ := entity.NewVisit(visitID, int64(15), sessionID, clientID, types.Hash{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(uuid.New().Generate(), clientID)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestRepository_Verify_EmptyClientID(t *testing.T) {
	_, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	visit, _ := entity.NewVisit(visitID, int64(15), sessionID, clientID, types.Hash{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify(uuid.New().Generate(), "")
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestRepository_Verify_EmptySessionID(t *testing.T) {
	_, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "test_FindClientID"

	visit, _ := entity.NewVisit(visitID, int64(15), sessionID, clientID, types.Hash{})

	_ = repository.Insert(visit)

	ok, err := repository.Verify([16]byte{}, clientID)
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestRepository_Insert(t *testing.T) {
	client, repository := Repository_CreateRepository()
	visitID := uuid.New().Generate()
	sessionID := uuid.New().Generate()
	clientID := "clientID"
	fields := types.Hash{"data": "here"}
	timestamp := int64(15)

	indexName := repository.indexName()
	typeName := repository.typeName

	visit, _ := entity.NewVisit(visitID, timestamp, sessionID, clientID, fields)

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
	assert.Equal(t, visit.Fields(), foundVisit.Fields())
}

func TestRepository_Insert_Nil(t *testing.T) {
	_, repository := Repository_CreateRepository()

	err := repository.Insert(nil)

	assert.NotNil(t, err)
}

func Repository_CreateRepository() (*driver.Client, *Repository) {
	client, _ := driver.NewClient()
	repository, _ := NewRepository(client, uuid.New())
	repository.indexPrefix = "test-tracker-"
	repository.RefreshAfterInsert = true

	_, _ = client.DeleteIndex(repository.indexName()).Do()
	_, _ = client.CreateIndex(repository.indexName()).Do()

	return client, repository
}
