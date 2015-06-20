package elastic

import (
	"encoding/json"
	"errors"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
)

type EventRepository struct {
	indexName string
	typeName  string
	client    *driver.Client
	uuid      dao.UUIDProviderInterface
}

func NewEventRepository(client *driver.Client, uuid dao.UUIDProviderInterface) (*EventRepository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	return &EventRepository{indexName: "tracker", typeName: "event", client: client, uuid: uuid}, nil
}

func (repository *EventRepository) FindAll() ([]*entities.Event, error) {
	return repository.find(nil, 0, 0)
}

func (repository *EventRepository) FindAllByVisit(visit *entities.Visit) ([]*entities.Event, error) {
	data := visit.Data()

	outer := driver.NewBoolFilter().Must(driver.NewTermFilter("enabled", true))

	rangeOuter := driver.NewBoolFilter().Should(driver.NewMissingFilter("filter.key"))

	keys := make([]string, len(data))

	var i uint
	for key, value := range data {
		inner := driver.NewBoolFilter().
			Must(driver.NewTermFilter("filter.key", key)).
			MustNot(driver.NewTermFilter("filter.value", value))

		nested := driver.NewNestedFilter("filterList").Filter(inner)

		outer.MustNot(nested)

		keys[i] = key
		i++
	}

	if i > 0 {
		rangeOuter.Should(driver.NewTermsFilter("filter.key", keys))
	}

	outer.Must(rangeOuter)

	return repository.find(&outer, 0, 0)
}

func (repository *EventRepository) FindByID(id [16]byte) (*entities.Event, error) {
	termQuery := driver.NewTermQuery("_id", repository.uuid.ToString(id))

	result, err := repository.find(&termQuery, 0, 1)

	if (err != nil) || (len(result) == 0) {
		return nil, err
	}

	return result[0], nil
}

func (repository *EventRepository) Insert(event *entities.Event) (err error) {
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

func (repository *EventRepository) Update(event *entities.Event) (err error) {
	if event == nil {
		return errors.New("Empty event is not allowed")
	}

	eventData, err := repository.eventToByte(event)

	if err != nil {
		return err
	}

	_, err = repository.client.Update().
		Index(repository.indexName).
		Type(repository.typeName).
		Id(repository.uuid.ToString(event.EventID())).
		Doc(string(eventData)).
		Refresh(true).
		Do()

	return err
}

func (repository *EventRepository) find(term driver.Query, limit, offset uint) ([]*entities.Event, error) {
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
		return []*entities.Event{}, err
	}

	result := make([]*entities.Event, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		event, err := repository.byteToEvent(*hit.Source)

		if err != nil {
			return []*entities.Event{}, err
		}

		result[i] = event
	}

	return result, nil
}

func (repository *EventRepository) eventToByte(event *entities.Event) ([]byte, error) {
	structEvent := eventStructEvent{
		EventID: repository.uuid.ToString(event.EventID()),
		Enabled: event.Enabled(),
	}

	dataFromEvent := event.Data()
	structEvent.DataList = make([]eventStructHash, len(dataFromEvent))

	var i uint
	for key, value := range dataFromEvent {
		structEvent.DataList[i] = eventStructHash{Key: key, Value: value}

		i++
	}

	filtersFromEvent := event.Filters()
	structEvent.FilterList = make([]eventStructHash, len(filtersFromEvent))

	i = 0
	for key, value := range dataFromEvent {
		structEvent.FilterList[i] = eventStructHash{Key: key, Value: value}

		i++
	}

	return json.Marshal(structEvent)
}

func (repository *EventRepository) byteToEvent(data []byte) (*entities.Event, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structEvent := new(eventStructEvent)

	err := json.Unmarshal(data, structEvent)
	if err != nil {
		return nil, err
	}

	dataList := make(map[string]string, len(structEvent.DataList))
	for _, value := range structEvent.DataList {
		dataList[value.Key] = value.Value
	}

	filtersList := make(map[string]string, len(structEvent.FilterList))
	for _, value := range structEvent.FilterList {
		filtersList[value.Key] = value.Value
	}

	return entities.NewEvent(repository.uuid.ToBytes(structEvent.EventID), structEvent.Enabled, dataList, filtersList)
}

type eventStructEvent struct {
	EventID    string            `json:"_id"`
	Enabled    bool              `json:"enabled"`
	DataList   []eventStructHash `json:"dataList"`
	FilterList []eventStructHash `json:"filterList"`
}

type eventStructHash struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
