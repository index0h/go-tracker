package elastic

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
)

type FlashRepository struct {
	RefreshAfterInsert bool
	indexPrefix        string
	typeName           string
	client             *driver.Client
	uuid               dao.UUIDProviderInterface
}

func NewFlashRepository(client *driver.Client, uuid dao.UUIDProviderInterface) (*FlashRepository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	return &FlashRepository{typeName: "event_flash", indexPrefix: "tracker-", client: client, uuid: uuid}, nil
}

func (repository *FlashRepository) find(term driver.Query, limit, offset uint) ([]*entities.Flash, error) {
	request := repository.client.
		Search().
		Index(repository.indexName()).
		Type(repository.typeName).
		Sort("@timestamp", false)

	if term != nil {
		request = request.Query(term)
	}

	if limit > 0 {
		request = request.From(int(limit))
	}

	if offset > 0 {
		request = request.Size(int(offset))
	}

	searchResult, err := request.Do()

	if (err != nil) || (searchResult.TotalHits() == 0) {
		return []*entities.Flash{}, err
	}

	result := make([]*entities.Flash, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		flash, err := repository.byteToFlash(*hit.Source)

		if err != nil {
			return []*entities.Flash{}, err
		}

		result[i] = flash
	}

	return result, nil
}

func (repository *FlashRepository) indexName() string {
	return repository.indexPrefix + time.Unix(time.Now().Unix(), 0).Format("2006-01")
}

func (repository *FlashRepository) flashToByte(flash *entities.Flash) ([]byte, error) {
	model := elasticFlash{
		FlashID:     repository.uuid.ToString(flash.VisitID()),
		Timestamp:   time.Unix(flash.Timestamp(), 0).Format("2006-01-02 15:04:05"),
		VisitID:     repository.uuid.ToString(flash.VisitID()),
		EventID:     repository.uuid.ToString(flash.EventID()),
		VisitFields: keyValFromHash(flash.VisitFields()),
		EventFields: keyValFromHash(flash.EventFields()),
	}

	return json.Marshal(model)
}

func (repository *FlashRepository) byteToFlash(data []byte) (*entities.Flash, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structFlash := new(elasticFlash)

	if err := json.Unmarshal(data, structFlash); err != nil {
		return nil, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", structFlash.Timestamp)
	if err != nil {
		return nil, err
	}

	return entities.NewFlashFromRaw(
		repository.uuid.ToBytes(structFlash.FlashID),
		repository.uuid.ToBytes(structFlash.VisitID),
		repository.uuid.ToBytes(structFlash.EventID),
		timestamp.Unix(),
		hashFromKeyVal(structFlash.VisitFields),
		hashFromKeyVal(structFlash.EventFields),
	)
}
