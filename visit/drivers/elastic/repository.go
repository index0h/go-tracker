package elastic

import (
	"errors"
	"encoding/json"
	"time"

	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
	elasticDriver "github.com/olivere/elastic"
)

const (
	indexPrefix = "track-visit-"
	indexSuffixLayout = "2006-01"
	typeName = "visit"
	timestampLayout = "2006-01-02 15:04:05"
	visitIDName = "_id"
	sessionIDName = "sessionID"
	clientIDName = "clientID"
	timestampName = "@timestamp"
)

type Repository struct {
	currentIndexName string
	client           *elasticDriver.Client
	uuid             uuid.Maker
}

func NewRepository(client *elasticDriver.Client, uuid uuid.Maker) *Repository {
	return &Repository{client: client, uuid: uuid}
}

// Find clientID by sessionID. If it's not present in cache - will try to find by nested repository and cache result
func (repository *Repository) FindClientID(sessionID uuid.UUID) (clientID string, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return clientID, errors.New("Empty sessionID is not allowed")
	}

	termQuery := elasticDriver.NewTermQuery(sessionIDName, repository.uuid.ToString(sessionID))

	visit, err := repository.searchOneVisitByTerm(&termQuery)

	if (err != nil) || (visit == nil) {
		return clientID, err
	}

	return visit.ClientID(), err
}

// Find sessionID by clientID. If it's not present in cache - will try to find by nested repository and cache result
func (repository *Repository) FindSessionID(clientID string) (sessionID uuid.UUID, err error) {
	if clientID == "" {
		return sessionID, errors.New("Empty clientID is not allowed")
	}

	termQuery := elasticDriver.NewTermQuery(clientIDName, clientID)

	visit, err := repository.searchOneVisitByTerm(&termQuery)

	if (err != nil) || (visit == nil) {
		return sessionID, err
	}

	return visit.SessionID(), err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
// If sessionID or clientID not found it'll run nested repository and cache result (if its ok)
func (repository *Repository) Verify(sessionID uuid.UUID, clientID string) (ok bool, err error) {
	if uuid.IsUUIDEmpty(sessionID) {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	boolFilter := elasticDriver.NewBoolFilter().
		Must(elasticDriver.NewTermFilter(sessionIDName, repository.uuid.ToString(sessionID))).
		MustNot(elasticDriver.NewTermFilter(clientIDName, clientID)).
		MustNot(elasticDriver.NewTermFilter(clientIDName, ""))

	var indexName string
	indexName, err = repository.indexName()

	if err != nil {
		return false, err
	}

	searchResult, err := repository.client.Search().
		Index(indexName).
		Query(&boolFilter).
		Sort(timestampName, false).
		From(0).Size(1).
		Do()

	if (err != nil) || (searchResult.TotalHits() > 0) {
		return false, err
	}

	return true, nil
}

// Save visit
func (repository *Repository) Insert(visit *entities.Visit) (err error) {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	visitData, err := repository.visitToByte(visit)

	if err != nil {
		return err
	}

	var indexName string
	indexName, err = repository.indexName()

	if err != nil {
		return err
	}

	_, err = repository.client.Index().
		Index(indexName).
		Type(typeName).
		Id(repository.uuid.ToString(visit.VisitID())).
		BodyString(string(visitData)).
		Do()

	return err
}

// Return current index name and check that it exists
func (repository *Repository) indexName() (result string, err error) {
	result = indexPrefix + time.Unix(time.Now().Unix(), 0).Format(indexSuffixLayout)

	if repository.currentIndexName != result {
		exists, err := repository.client.IndexExists(result).Do()

		if exists {
			return result, nil
		}

		if err != nil {
			return result, err
		}

		repository.currentIndexName = result
	}

	return result, nil
}

// Search one visit by filter term
func (repository *Repository) searchOneVisitByTerm(term *elasticDriver.TermQuery) (visit *entities.Visit, err error) {
	var indexName string
	indexName, err = repository.indexName()

	if err != nil {
		return visit, err
	}

	searchResult, err := repository.client.Search().
		Index(indexName).
		Query(term).
		Sort(timestampName, false).
		From(0).Size(1).
		Do()

	if (err != nil) || (searchResult.TotalHits() == 0) {
		return visit, err
	}

	rawVisit := new(mapVisit)

	json.Unmarshal(*searchResult.Hits.Hits[0].Source, rawVisit)

	return repository.rawToVisit(rawVisit)
}

// Convert visit to bytes
func (repository *Repository) visitToByte(visit *entities.Visit) ([]byte, error) {
	model := mapVisit{
		VisitID:     repository.uuid.ToString(visit.VisitID()),
		Timestamp:   time.Unix(visit.Timestamp(), 0).Format(timestampLayout),
		SessionID:   repository.uuid.ToString(visit.SessionID()),
		ClientID:    visit.ClientID(),
		WarningList: visit.Warnings(),
	}

	dataFromVisit := visit.Data()
	model.DataList = make([]mapDataList, len(dataFromVisit))

	var i uint
	for key, value := range dataFromVisit {
		model.DataList[i] = mapDataList{Key: key, Value: value}

		i++
	}

	return json.Marshal(model)
}

// Convert mapVisit to Visit instance
func (repository *Repository) rawToVisit(rawVisit *mapVisit) (visit *entities.Visit, err error) {
	timestamp, err := time.Parse(timestampLayout, rawVisit.Timestamp)

	if err != nil {
		return visit, err
	}

	dataList := make(map[string]string, len(rawVisit.DataList))
	for _, value := range rawVisit.DataList {
		dataList[value.Key] = value.Value
	}

	return entities.NewVisit(
		repository.uuid.ToBytes(rawVisit.VisitID),
		timestamp.Unix(),
		repository.uuid.ToBytes(rawVisit.SessionID),
		rawVisit.ClientID,
		dataList,
		rawVisit.WarningList,
	)
}

// Return visit index mapping
func (repository *Repository) indexBody() string {
	return `{
  "mapping":{
    "visit":{
      "properties":{
        "_id":{"index":"not_analyzed", "stored":true, "type":"string"},
        "@timestamp":{"format":"YYYY-MM-DD HH:mm:ss", "type":"date"},
        "clientId":{"index":"not_analyzed", "type":"string"},
        "dataList":{
          "include_in_parent":true,
          "type":"nested",
          "properties":{
            "key":{"index":"not_analyzed", "type":"string"},
            "value":{"index":"not_analyzed", "type":"string"}
          }
        },
        "sessionId":{"index":"not_analyzed", "type":"string"},
        "warnings":{"index":"not_analyzed", "type":"string"}
      }
    }
  }
}`
}

type mapVisit struct {
	VisitID     string        `json:"_id"`
	Timestamp   string        `json:"@timestamp"`
	SessionID   string        `json:"sessionId"`
	ClientID    string        `json:"clientId"`
	DataList    []mapDataList `json:"dataList"`
	WarningList []string      `json:"warningList"`
}

type mapDataList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
