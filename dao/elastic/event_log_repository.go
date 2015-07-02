package elastic

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
)

type EventLogRepository struct {
	RefreshAfterInsert bool
	indexPrefix        string
	typeName           string
	client             *driver.Client
	uuid               dao.UUIDProviderInterface
}

func NewEventLogRepository(client *driver.Client, uuid dao.UUIDProviderInterface) (*EventLogRepository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	return &EventLogRepository{typeName: "event_log", indexPrefix: "tracker-", client: client, uuid: uuid}, nil
}

func (repository *EventLogRepository) find(term driver.Query, limit, offset uint) ([]*entities.EventLog, error) {
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
		return []*entities.EventLog{}, err
	}

	result := make([]*entities.EventLog, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		eventLog, err := repository.byteToEventLog(*hit.Source)

		if err != nil {
			return []*entities.EventLog{}, err
		}

		result[i] = eventLog
	}

	return result, nil
}

func (repository *EventLogRepository) indexName() string {
	return repository.indexPrefix + time.Unix(time.Now().Unix(), 0).Format("2006-01")
}

func (repository *EventLogRepository) eventLogToByte(eventLog *entities.EventLog) ([]byte, error) {
	model := eventLogStructEventLog{
		EventLogID: repository.uuid.ToString(eventLog.VisitID()),
		Timestamp:  time.Unix(eventLog.Timestamp(), 0).Format("2006-01-02 15:04:05"),
		VisitID:    repository.uuid.ToString(eventLog.VisitID()),
	}

	visitDataList := eventLog.VisitData()
	model.VisitDataList = make([]eventLogStructHash, len(visitDataList))

	var i uint
	for key, value := range visitDataList {
		model.VisitDataList[i] = eventLogStructHash{Key: key, Value: value}

		i++
	}

	i = 0
	eventsData := eventLog.EventsData()
	model.EventList = make([]eventLogStructEvent, len(eventsData))
	for eventID, dataList := range eventsData {
		event := eventLogStructEvent{
			EventID: repository.uuid.ToString(eventID),
			EventDataList: make([]eventLogStructHash, len(dataList)),
		}

		var j uint
		for key, value := range dataList {
			event.EventDataList[j] = eventLogStructHash{Key: key, Value: value}

			j++
		}

		model.EventList[i] = event
		i++
	}

	return json.Marshal(model)
}

func (repository *EventLogRepository) byteToEventLog(data []byte) (*entities.EventLog, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structEventLog := new(eventLogStructEventLog)

	if err := json.Unmarshal(data, structEventLog); err != nil {
		return nil, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", structEventLog.Timestamp)
	if err != nil {
		return nil, err
	}

	visitDataList := make(map[string]string, len(structEventLog.VisitDataList))
	for _, value := range structEventLog.VisitDataList {
		visitDataList[value.Key] = value.Value
	}

	eventsData := make(map[[16]byte]map[string]string, len(structEventLog.EventList))
	for _, value := range structEventLog.EventList {
		eventDataList := make(map[string]string, len(value.EventDataList))
		for _, value := range structEventLog.VisitDataList {
			eventDataList[value.Key] = value.Value
		}

		eventsData[repository.uuid.ToBytes(value.EventID)] = eventDataList
	}

	return entities.NewEventLogFromRaw(
		repository.uuid.ToBytes(structEventLog.EventLogID),
		timestamp.Unix(),
		repository.uuid.ToBytes(structEventLog.VisitID),
		visitDataList,
		eventsData,
	)
}

type eventLogStructEventLog struct {
	EventLogID    string                `json:"_id"`
	Timestamp     string                `json:"@timestamp"`
	VisitID       string                `json:"visitId"`
	VisitDataList []eventLogStructHash  `json:"visitDataList"`
	EventList     []eventLogStructEvent `json:"eventList"`
}

type eventLogStructHash struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type eventLogStructEvent struct {
	EventID       string               `json:"eventId"`
	EventDataList []eventLogStructHash `json:"eventDataList"`
}
