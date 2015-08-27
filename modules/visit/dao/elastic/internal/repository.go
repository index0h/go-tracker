package internal

import (
	"encoding/json"
	"errors"
	"time"

	shareElastic "github.com/index0h/go-tracker/share/elastic"
	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/modules/visit/entity"
)

type Repository struct {
	uuid share.UUIDProviderInterface
}

func (repository *Repository) Refresh(scenario int) bool {
	switch scenario {
	case shareElastic.BulkScenario:
		return true
	case shareElastic.InsertScenario:
		return false
	default:
		panic(errors.New("invalid scenario"))
	}
}

func (repository *Repository) GetIndexName(scenario int) string {
	switch scenario {
	case shareElastic.BulkScenario:
		fallthrough
	case shareElastic.InsertScenario:
		return "tracker-visit-" + time.Unix(time.Now().Unix(), 0).Format("2006-01")
	case shareElastic.FindScenario:
		return "tracker-visit-*"
	default:
		panic(errors.New("invalid scenario"))
	}
}

func (repository *Repository) GetTypeName() string {
	return "visit"
}

func (repository *Repository) Marshal(visit interface{}) ([]byte, error) {
	entity := visit.(entity.Visit)

	model := Visit{
		VisitID:   repository.uuid.ToString(entity.VisitID()),
		Timestamp: time.Unix(entity.Timestamp(), 0).Format("2006-01-02 15:04:05"),
		SessionID: repository.uuid.ToString(entity.SessionID()),
		ClientID:  entity.ClientID(),
		Fields:    keyValFromHash(entity.Fields()),
	}

	return json.Marshal(model)
}

func (repository *Repository) Unmarshal(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data is not allowed")
	}

	visit := new(Visit)

	err := json.Unmarshal(data, visit)
	if err != nil {
		return nil, err
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", visit.Timestamp)

	if err != nil {
		return nil, err
	}

	return entity.NewVisit(
		repository.uuid.ToBytes(visit.VisitID),
		timestamp.Unix(),
		repository.uuid.ToBytes(visit.SessionID),
		visit.ClientID,
		hashFromKeyVal(visit.Fields),
	)
}
