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
	indexPrefix        string
	typeName           string
	client             *driver.Client
	uuid               dao.UUIDProviderInterface
}

func NewVisitRepository(client *driver.Client, uuid dao.UUIDProviderInterface) (*VisitRepository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	return &VisitRepository{typeName: "visit", indexPrefix: "tracker-", client: client, uuid: uuid}, nil
}

func (repository *VisitRepository) FindByID(visitID [16]byte) (*entities.Visit, error) {
	if visitID == [16]byte{} {
		return nil, errors.New("Empty visitID is not allowed")
	}

	search, err := repository.client.Get().
		Index("twitter").
		Type("tweet").
		Id("1").
		Do()

	if err != nil {
		return nil, err
	}

	if !search.Found {
		return nil, nil
	}

	return repository.byteToVisit(*search.Source)
}

func (repository *VisitRepository) FindAll(limit int64, offset int64) ([]*entities.Visit, error) {
	return repository.find(nil, uint(limit), uint(offset))
}

func (repository *VisitRepository) FindAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) (result []*entities.Visit, err error) {
	if sessionID == [16]byte{} {
		return result, errors.New("Empty sessionID is not allowed")
	}

	termQuery := driver.NewTermQuery("sessionId", repository.uuid.ToString(sessionID))

	return repository.find(&termQuery, uint(limit), uint(offset))
}

func (repository *VisitRepository) FindAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) (result []*entities.Visit, err error) {
	if clientID == "" {
		return result, errors.New("Empty clientID is not allowed")
	}

	termQuery := driver.NewTermQuery("clientID", clientID)

	return repository.find(&termQuery, uint(limit), uint(offset))
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

	if repository.RefreshAfterInsert {
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
	model := elasticVisit{
		VisitID:   repository.uuid.ToString(visit.VisitID()),
		Timestamp: time.Unix(visit.Timestamp(), 0).Format("2006-01-02 15:04:05"),
		SessionID: repository.uuid.ToString(visit.SessionID()),
		ClientID:  visit.ClientID(),
		Fields:    keyValFromHash(visit.Fields()),
	}

	return json.Marshal(model)
}

func (repository *VisitRepository) byteToVisit(data []byte) (*entities.Visit, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	structVisit := new(elasticVisit)

	err := json.Unmarshal(data, structVisit)
	if err != nil {
		return nil, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", structVisit.Timestamp)

	if err != nil {
		return nil, err
	}

	return entities.NewVisit(
		repository.uuid.ToBytes(structVisit.VisitID),
		timestamp.Unix(),
		repository.uuid.ToBytes(structVisit.SessionID),
		structVisit.ClientID,
		hashFromKeyVal(structVisit.Fields),
	)
}
