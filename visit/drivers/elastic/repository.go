package elastic

import (
	"errors"
	"encoding/json"
	"time"

	"github.com/index0h/go-tracker/uuid"
	"github.com/index0h/go-tracker/visit/entities"
	elasticDriver "github.com/olivere/elastic"
	"fmt"
)

const (
	indexPrefix = "nya-"
	indexSuffixLayout = "2006-01"
	typeName = "visit"
	timestampLayout = "2006-01-02 15:04:05"
	visitIDName = "_id"
	sessionIDName = "sessionID"
	clientIDName = "clientID"
	timestampName = "@timestamp"
)

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

type Repository struct {
	currentIndexName string
	client           *elasticDriver.Client
	uuid             uuid.Maker
}

func NewRepository(client *elasticDriver.Client, uuid uuid.Maker) *Repository {
	return &Repository{client: client, uuid: uuid}
}

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

func (repository *Repository) FindSessionID(clientID string) (sessionID uuid.UUID, err error) {
	if clientID == "" {
		return sessionID, errors.New("Empty clientID is not allowed")
	}

	termQuery := elasticDriver.NewTermQuery(clientIDName, clientID)

	visit, err := repository.searchOneVisitByTerm(&termQuery)

	fmt.Println(visit, "NYAA")

	if (err != nil) || (visit == nil) {
		return sessionID, err
	}

	return visit.SessionID(), err
}

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

	searchResult, err := repository.client.Search().
		Index(repository.indexName()).
		Query(&boolFilter).
		Sort(timestampName, false).
		From(0).Size(1).
		Do()

	if (err != nil) || (searchResult.TotalHits() > 0) {
		return false, err
	}

	return true, nil
}

func (repository *Repository) Insert(visit *entities.Visit) (err error) {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	repository.client.CreateIndex(repository.indexName()).Body(repository.indexBody()).Do()

	visitData, err := repository.visitToByte(visit)

	if err != nil {
		return err
	}

	_, err = repository.client.Index().
		Index(repository.indexName()).
		Type(typeName).
		Id(repository.uuid.ToString(visit.VisitID())).
		BodyString(string(visitData)).
		Do()

	return err
}

func (repository *Repository) searchOneVisitByTerm(term *elasticDriver.TermQuery) (visit *entities.Visit, err error) {
	searchResult, err := repository.client.Search().
		Index(repository.indexName()).
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

func (repository *Repository) indexName() string {
	return indexPrefix + time.Unix(time.Now().Unix(), 0).Format(indexSuffixLayout)
}

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
