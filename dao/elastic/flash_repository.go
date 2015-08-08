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

func (repository *FlashRepository) FindByID(flashID [16]byte) (result *entities.Flash, err error) {
	if flashID == [16]byte{} {
		return nil, errors.New("Empty flashID is not allowed")
	}

	search, err := repository.client.Get().
		Index(repository.indexName()).
		Type(repository.typeName).
		Id(repository.uuid.ToString(flashID)).
		Do()

	if err != nil {
		return nil, err
	}

	if !search.Found {
		return nil, nil
	}

	return repository.byteToFlash(*search.Source)
}

func (repository *FlashRepository) FindAll(limit int64, offset int64) (result []*entities.Flash, err error) {
	return repository.find(nil, uint(limit), uint(offset))
}

func (repository *FlashRepository) FindAllByVisitID(visitID [16]byte) (result []*entities.Flash, err error) {
	if visitID == [16]byte{} {
		return result, errors.New("Empty visitID is not allowed")
	}

	termQuery := driver.NewTermQuery("visitId", repository.uuid.ToString(visitID))

	return repository.find(&termQuery, 0, 0)
}

func (repository *FlashRepository) FindAllByEventID(
	eventID [16]byte,
	limit int64,
	offset int64,
) (result []*entities.Flash, err error) {
	if eventID == [16]byte{} {
		return result, errors.New("Empty eventID is not allowed")
	}

	termQuery := driver.NewTermQuery("eventID", repository.uuid.ToString(eventID))

	return repository.find(&termQuery, uint(limit), uint(offset))
}

func (repository *FlashRepository) Insert(flash *entities.Flash) (err error) {
	if flash == nil {
		return errors.New("Empty flash is not allowed")
	}

	flashData, err := repository.flashToByte(flash)

	if err != nil {
		return err
	}

	request := repository.client.
		Index().
		Index(repository.indexName()).
		Type(repository.typeName).
		Id(repository.uuid.ToString(flash.FlashID())).
		BodyString(string(flashData))

	if repository.RefreshAfterInsert {
		request.Refresh(true)
	}

	_, err = request.Do()

	return err
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
