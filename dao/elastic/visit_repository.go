package elastic

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	driver "github.com/olivere/elastic"
)

type VisitRepository struct {
	RefreshAfterInsert bool
	indexPrefix string
	typeName string
	client   *driver.Client
	uuid     dao.UUIDProviderInterface
}

func NewVisitRepository(client *driver.Client, uuid dao.UUIDProviderInterface) *VisitRepository {
	return &VisitRepository{typeName: "visit", indexPrefix: "tracker-", client: client, uuid: uuid}
}

// Find clientID by sessionID. If it's not present in cache - will try to find by nested repository and cache result
func (repository *VisitRepository) FindClientID(sessionID [16]byte) (clientID string, err error) {
	if sessionID == [16]byte{} {
		return clientID, errors.New("Empty sessionID is not allowed")
	}

	termQuery := driver.NewTermQuery("sessionId", repository.uuid.ToString(sessionID))

	visit, err := repository.find(&termQuery, 0, 1)

	if (err != nil) || (len(visit) == 0) {
		return clientID, err
	}

	return visit[0].ClientID(), err
}

// Verify method MUST check that sessionID is not registered by another not empty clientID
func (repository *VisitRepository) Verify(sessionID [16]byte, clientID string) (ok bool, err error) {
	if sessionID == [16]byte{} {
		return false, errors.New("Empty sessioID is not allowed")
	}

	if clientID == "" {
		return false, errors.New("Empty clientID is not allowed")
	}

	boolFilter := driver.NewBoolFilter().
		Must(driver.NewTermFilter("sessionId", repository.uuid.ToString(sessionID))).
		MustNot(driver.NewTermFilter("clientId", clientID)).
		MustNot(driver.NewTermFilter("clientId", ""))

	visits, err := repository.find(&boolFilter, 0, 1)

	if (err != nil) || (len(visits) > 0) {
		return false, err
	}

	return true, nil
}

// Save visit
func (repository *VisitRepository) Insert(visit *entities.Visit) (err error) {
	if visit == nil {
		return errors.New("Empty visit is not allowed")
	}

	visitData, err := repository.visitToByte(visit)

	if err != nil {
		return err
	}

	request := repository.client.
		Index().
		Index(repository.indexName()).
		Type(repository.typeName).
		Id(repository.uuid.ToString(visit.VisitID())).
		BodyString(string(visitData))

	if (repository.RefreshAfterInsert) {
		request.Refresh(true)
	}

	_, err = request.Do()

	return err
}

func (repository *VisitRepository) find(term driver.Query, limit, offset uint) ([]*entities.Visit, error) {
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
		return []*entities.Visit{}, err
	}

	result := make([]*entities.Visit, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		visit, err := repository.byteToVisit(*hit.Source)

		if err != nil {
			return []*entities.Visit{}, err
		}

		result[i] = visit
	}

	return result, nil
}

// Return current index name and check that it exists
func (repository *VisitRepository) indexName() string {
	return repository.indexPrefix + time.Unix(time.Now().Unix(), 0).Format("2006-01")
}

// Convert visit to bytes
func (repository *VisitRepository) visitToByte(visit *entities.Visit) ([]byte, error) {
	model := visitStructVisit{
		VisitID:     repository.uuid.ToString(visit.VisitID()),
		Timestamp:   time.Unix(visit.Timestamp(), 0).Format("2006-01-02 15:04:05"),
		SessionID:   repository.uuid.ToString(visit.SessionID()),
		ClientID:    visit.ClientID(),
		WarningList: visit.Warnings(),
	}

	dataFromVisit := visit.Data()
	model.DataList = make([]visitStructHash, len(dataFromVisit))

	var i uint
	for key, value := range dataFromVisit {
		model.DataList[i] = visitStructHash{Key: key, Value: value}

		i++
	}

	return json.Marshal(model)
}

func (repository *VisitRepository) byteToVisit(data []byte) (*entities.Visit, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structVisit := new(visitStructVisit)

	err := json.Unmarshal(data, structVisit)
	if err != nil {
		return nil, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", structVisit.Timestamp)

	if err != nil {
		return nil, err
	}

	dataList := make(map[string]string, len(structVisit.DataList))
	for _, value := range structVisit.DataList {
		dataList[value.Key] = value.Value
	}

	return entities.NewVisit(
		repository.uuid.ToBytes(structVisit.VisitID),
		timestamp.Unix(),
		repository.uuid.ToBytes(structVisit.SessionID),
		structVisit.ClientID,
		dataList,
		structVisit.WarningList,
	)
}

type visitStructVisit struct {
	VisitID     string            `json:"_id"`
	Timestamp   string            `json:"@timestamp"`
	SessionID   string            `json:"sessionId"`
	ClientID    string            `json:"clientId"`
	DataList    []visitStructHash `json:"dataList"`
	WarningList []string          `json:"warningList"`
}

type visitStructHash struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
