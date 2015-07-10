package elastic

import "github.com/index0h/go-tracker/entities"

type elasticVisit struct {
	VisitID   string   `json:"_id"`
	Timestamp string   `json:"@timestamp"`
	SessionID string   `json:"sessionId"`
	ClientID  string   `json:"clientId"`
	Fields    []keyVal `json:"fields"`
}

type elasticEvent struct {
	EventID string   `json:"_id"`
	Enabled bool     `json:"enabled"`
	Fields  []keyVal `json:"fields"`
	Filters []keyVal `json:"filters"`
}

type elasticFlash struct {
	FlashID     string   `json:"_id"`
	VisitID     string   `json:"visitId"`
	EventID     string   `json:"eventId"`
	Timestamp   string   `json:"@timestamp"`
	VisitFields []keyVal `json:"visitFields"`
	EventFields []keyVal `json:"eventFields"`
}

type keyVal struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func keyValFromHash(data entities.Hash) []keyVal {
	result := make([]keyVal, len(data))

	var i uint
	for key, value := range data {
		result[i] = keyVal{Key: key, Value: value}

		i++
	}

	return result
}

func hashFromKeyVal(data []keyVal) entities.Hash {
	result := make(entities.Hash, len(data))

	for _, element := range data {
		result[element.Key] = element.Value
	}

	return result
}
