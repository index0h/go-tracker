package visitDriver

import (
	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entity"
	"github.com/olivere/elastic"
	"time"
)

type Repository struct {
	index *IndexService
	client *elastic.Client
	uuid   uuid.UuidMaker
}

func NewRepository(client *elastic.Client, uuid uuid.UuidMaker) *Repository {
	return &Repository{index: NewIndexService(uuid), client: client, uuid: uuid}
}

func (repository *Repository) Verify(sessionID uuid.Uuid, clientID string) (ok bool, err error) {
	return true, nil
}

func (repository *Repository) FindClientID(sessionID uuid.Uuid) (clientID string, err error) {
	return "", nil
}

func (repository *Repository) FindSessionID(clientID string) (sessionID uuid.Uuid, err error) {
	return uuid.Uuid{}, nil
}

func (repository *Repository) Insert(visit *entity.Visit) (err error) {
	repository.client.CreateIndex(repository.index.Name(time.Now().Unix())).Body(repository.index.Body()).Do()

	visitID, visitJson, err := repository.index.Marshal(visit)
	if err != nil {
		return err
	}

	_, err = repository.client.Index().
	Index(repository.index.Name(time.Now().Unix())).
	Type(repository.index.Type()).
	Id(visitID).
	BodyString(string(visitJson)).
	Do()

	return err
}