package elastic

import (
	"errors"
	"time"

	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
	elasticDriver "github.com/olivere/elastic"
)

type Repository struct {
	index  *IndexService
	client *elasticDriver.Client
	uuid   uuid.Maker
}

func NewRepository(client *elasticDriver.Client, uuid uuid.Maker) *Repository {
	return &Repository{index: NewIndexService(uuid), client: client, uuid: uuid}
}

func (repository *Repository) FindClientID(sessionID uuid.UUID) (clientID string, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return clientID, errors.New("Empty sessioID is not allowed")
	}
	// TODO: implement

	return clientID, err
}

func (repository *Repository) FindSessionID(clientID string) (sessionID uuid.UUID, err error) {
	if clientID == "" {
		return sessionID, errors.New("Empty clientID is not allowed")
	}
	// TODO: implement

	return sessionID, err
}

func (repository *Repository) Verify(sessionID uuid.UUID, clientID string) (ok bool, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}
	// TODO: implement

	return true, nil
}

// TODO: move index creation to another method
func (repository *Repository) Insert(visit *entities.Visit) (err error) {
	indexName := repository.index.Name(time.Now().Unix())

	repository.client.CreateIndex(indexName).Body(repository.index.Body()).Do()

	visitID, visitJSON, err := repository.index.Marshal(visit)
	if err != nil {
		return err
	}

	_, err = repository.client.Index().
		Index(indexName).
		Type(repository.index.Type()).
		Id(visitID).
		BodyString(visitJSON).
		Do()

	return err
}
