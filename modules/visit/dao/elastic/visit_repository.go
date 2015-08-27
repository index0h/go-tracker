package elastic

import (
	"errors"
	"time"

	"github.com/index0h/go-tracker/modules/visit/dao/elastic/internal"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share/elastic"
	"github.com/index0h/go-tracker/share"
	driver "github.com/olivere/elastic"
)

type VisitRepository struct {
	shareRepository *elastic.Repository
	uuid               share.UUIDProviderInterface
}

func NewVisitRepository(client *driver.Client, uuid share.UUIDProviderInterface) (*VisitRepository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	if uuid == nil {
		return nil, errors.New("uuid must be not nil")
	}

	shareRepository, _ := elastic.NewRepository(client, &internal.Repository{uuid: uuid}, 10 * time.Second)

	return &VisitRepository{shareRepository: shareRepository, uuid: uuid}, nil
}

func (repository *VisitRepository) FindByID(visitID [16]byte) (*entity.Visit, error) {
	if visitID == [16]byte{} {
		return nil, errors.New("Empty visitID is not allowed")
	}

	search, err := repository.client.Get().
		Index(repository.indexName()).
		Type(repository.typeName).
		Id(repository.uuid.ToString(visitID)).
		Do()

	if err != nil {
		return nil, err
	}

	if !search.Found {
		return nil, nil
	}

	return repository.byteToVisit(*search.Source)
}

func (repository *VisitRepository) FindAll(limit int64, offset int64) ([]*entity.Visit, error) {
	return repository.find(nil, uint(limit), uint(offset))
}

func (repository *VisitRepository) FindAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) (result []*entity.Visit, err error) {
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
) (result []*entity.Visit, err error) {
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
func (repository *VisitRepository) Insert(visit *entity.Visit) (err error) {
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

func (repository *VisitRepository) find(term driver.Query, limit, offset uint) ([]*entity.Visit, error) {
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
		return []*entity.Visit{}, err
	}

	result := make([]*entity.Visit, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		visit, err := repository.byteToVisit(*hit.Source)

		if err != nil {
			return []*entity.Visit{}, err
		}

		result[i] = visit
	}

	return result, nil
}
