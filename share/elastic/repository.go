package elastic

import (
	"errors"
	"sync"
	"time"

	driver "github.com/olivere/elastic"
)

const (
	FindScenario = iota
	InsertScenario
	DeleteScenario
	BulkScenario
)

type Repository struct {
	sync.Mutex
	client      *driver.Client
	bulkTime    int
	bulkService *driver.BulkService
	internal    InternalRepositoryInterface
}

func NewRepository(client *driver.Client, internal InternalRepositoryInterface, bulkTime uint) (*Repository, error) {
	if client == nil {
		return nil, errors.New("client must be not nil")
	}

	result := &Repository{client: client, internal: internal, bulkTime: int(bulkTime), bulkService: client.Bulk()}
	result.bulkLoop()

	return result, nil
}

func (repository *Repository) FindByID(id string) (interface{}, error) {
	if id == "" {
		return nil, errors.New("Empty ID is not allowed")
	}

	search, err := repository.
		client.
		Get().
		Index(repository.internal.GetIndexName(FindScenario)).
		Type(repository.internal.GetTypeName()).
		Id(id).
		Do()

	if err != nil {
		return nil, err
	}

	if !search.Found {
		return nil, nil
	}

	return repository.internal.Unmarshal(*search.Source)
}

func (repository *Repository) FindAll(limit int64, offset int64) ([]interface{}, error) {
	return repository.FindAllByQuery(nil, uint(limit), uint(offset))
}

func (repository *Repository) FindAllByQuery(term driver.Query, limit, offset uint) (result []interface{}, err error) {
	request := repository.
		client.
		Search().
		Index(repository.internal.GetIndexName(FindScenario)).
		Type(repository.internal.GetTypeName())

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
		return result, err
	}

	result = make([]interface{}, searchResult.TotalHits())

	for i, hit := range searchResult.Hits.Hits {
		visit, err := repository.internal.Unmarshal(*hit.Source)

		if err != nil {
			return []interface{}{}, err
		}

		result[i] = visit
	}

	return result, nil
}

func (repository *Repository) Insert(entity interface{}) (err error) {
	if entity == nil {
		return errors.New("Empty entity is not allowed")
	}

	id, err := repository.internal.GetEntityID(entity)
	if err != nil {
		return err
	}

	data, err := repository.internal.Marshal(entity)

	if err != nil {
		return err
	}

	request := repository.
		client.
		Index().
		Index(repository.internal.GetIndexName(InsertScenario)).
		Type(repository.internal.GetTypeName()).
		Id(id).
		BodyString(data).
		Refresh(repository.internal.Refresh(InsertScenario))

	_, err = request.Do()

	return err
}

func (repository *Repository) InsertAsync(entity interface{}) (err error) {
	if entity == nil {
		return errors.New("Empty entity is not allowed")
	}

	id, err := repository.internal.GetEntityID(entity)
	if err != nil {
		return err
	}

	data, err := repository.internal.Marshal(entity)

	if err != nil {
		return err
	}

	request := driver.
		NewBulkIndexRequest().
		Index().
		Index(repository.internal.GetIndexName(InsertScenario)).
		Type(repository.internal.GetTypeName()).
		Id(id).
		Doc(data)

	repository.Lock()

	repository.bulkService.Add(request)

	repository.Unlock()

	return err
}

func (repository *Repository) Delete(id string) (err error) {
	if id == "" {
		return errors.New("Empty id is not allowed")
	}

	request := repository.
		client.
		Delete().
		Index(repository.internal.GetIndexName(DeleteScenario)).
		Type(repository.internal.GetTypeName()).
		Id(id).
		Refresh(repository.internal.Refresh(DeleteScenario))

	_, err = request.Do()

	return err
}

func (repository *Repository) DeleteAsync(id string) (err error) {
	if id == "" {
		return errors.New("Empty id is not allowed")
	}

	request := driver.
		NewBulkDeleteRequest().
		Index().
		Index(repository.internal.GetIndexName(DeleteScenario)).
		Type(repository.internal.GetTypeName()).
		Id(id)

	repository.Lock()

	repository.bulkService.Add(request)

	repository.Unlock()

	return err
}

func (repository *Repository) bulkLoop() {
	runner := func(bulk *driver.BulkService) {
		bulk.Refresh(repository.internal.Refresh(BulkScenario)).Do()
	}

	for {
		time.Sleep(repository.bulkTime)

		repository.Lock()

		if repository.bulkService.NumberOfActions() == 0 {
			repository.Unlock()

			continue
		}

		go runner(repository.bulkService)

		repository.bulkService = repository.client.Bulk()

		repository.Unlock()
	}
}
