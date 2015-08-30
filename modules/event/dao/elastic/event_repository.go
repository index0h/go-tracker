package elastic

import (
	"encoding/json"
	"errors"

	"github.com/index0h/go-tracker/modules/event/entity"
	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/share/types"
	driver "github.com/olivere/elastic"
)

type Repository struct {
	indexName string
	typeName  string
	client    *driver.Client
	uuid      share.UUIDProviderInterface
}

func NewRepository(client *driver.Client, uuid share.UUIDProviderInterface) (*Repository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	return &Repository{indexName: "tracker", typeName: "event", client: client, uuid: uuid}, nil
}

func (repository *Repository) FindAll(limit int64, offset int64) ([]*entity.Event, error) {
	return repository.find(nil, uint(limit), uint(offset))
}

func (repository *Repository) FindAllByFields(fields types.Hash) ([]*entity.Event, error) {
	if fields == nil {
		return []*entity.Event{}, errors.New("Empty fields not allowed")
	}

	outer := driver.NewBoolFilter().Must(driver.NewTermFilter("enabled", true))

	keys := make([]string, len(fields))

	var i uint
	for key, value := range fields {
		inner := driver.NewBoolFilter().
			Must(driver.NewTermFilter("filters.key", key)).
			MustNot(driver.NewTermFilter("filters.value", value))

		nested := driver.NewNestedFilter("filters").Filter(inner)

		outer = outer.MustNot(nested)

		keys[i] = key
		i++
	}

	postFilter := driver.NewBoolFilter().Should(driver.NewMissingFilter("filters.key"))

	if i > 0 {
		postFilter = postFilter.Should(driver.NewBoolFilter().Must(driver.NewTermsFilter("filters.key", keys)))
	}

	searchResult, err := repository.client.
		Search().
		Index(repository.indexName).
		Type(repository.typeName).
		Query(outer).
		PostFilter(postFilter).
		Do()

	if (err != nil) || (searchResult.TotalHits() == 0) {
		return []*entity.Event{}, err
	}

	result := make([]*entity.Event, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		event, err := repository.byteToEvent(*hit.Source)

		if err != nil {
			return []*entity.Event{}, err
		}

		result[i] = event
	}

	return result, nil
}

func (repository *Repository) FindByID(eventID types.UUID) (*entity.Event, error) {
	if eventID.IsEmpty() {
		return nil, errors.New("Empty id is not allowed")
	}

	termQuery := driver.NewTermQuery("_id", repository.uuid.ToString(eventID))

	result, err := repository.find(&termQuery, 0, 1)

	if (err != nil) || (len(result) == 0) {
		return nil, err
	}

	return result[0], nil
}

func (repository *Repository) Insert(event *entity.Event) (err error) {
	if event == nil {
		return errors.New("Empty event is not allowed")
	}

	eventData, err := repository.eventToByte(event)

	if err != nil {
		return err
	}

	_, err = repository.client.Index().
		Index(repository.indexName).
		Type(repository.typeName).
		Id(repository.uuid.ToString(event.EventID())).
		BodyString(string(eventData)).
		Refresh(true).
		Do()

	return err
}

func (repository *Repository) Update(event *entity.Event) (err error) {
	if event == nil {
		return errors.New("Empty event is not allowed")
	}

	return repository.Insert(event)
}

func (repository *Repository) find(term driver.Query, limit, offset uint) ([]*entity.Event, error) {
	request := repository.client.
		Search().
		Index(repository.indexName).
		Type(repository.typeName)

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
		return []*entity.Event{}, err
	}

	result := make([]*entity.Event, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		event, err := repository.byteToEvent(*hit.Source)

		if err != nil {
			return []*entity.Event{}, err
		}

		result[i] = event
	}

	return result, nil
}

func (repository *Repository) eventToByte(event *entity.Event) ([]byte, error) {
	structEvent := elasticEvent{
		EventID: repository.uuid.ToString(event.EventID()),
		Enabled: event.Enabled(),
		Fields:  keyValFromHash(event.Fields()),
		Filters: keyValFromHash(event.Filters()),
	}

	return json.Marshal(structEvent)
}

func (repository *Repository) byteToEvent(data []byte) (*entity.Event, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structEvent := new(elasticEvent)

	if err := json.Unmarshal(data, structEvent); err != nil {
		return nil, err
	}

	return entity.NewEvent(
		repository.uuid.FromString(structEvent.EventID),
		structEvent.Enabled,
		hashFromKeyVal(structEvent.Fields),
		hashFromKeyVal(structEvent.Filters),
	)
}
