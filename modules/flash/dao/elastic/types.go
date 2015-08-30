package elastic

import (
	"github.com/index0h/go-tracker/share/types"
)

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

func keyValFromHash(data types.Hash) []keyVal {
	result := make([]keyVal, len(data))

	var i uint
	for key, value := range data {
		result[i] = keyVal{Key: key, Value: value}

		i++
	}

	return result
}

func hashFromKeyVal(data []keyVal) types.Hash {
	result := make(types.Hash, len(data))

	for _, element := range data {
		result[element.Key] = element.Value
	}

	return result
}
