package elastic

import "github.com/index0h/go-tracker/visit/entities"

type Indexer interface {
	//
	Name(timestamp int64) string

	//
	TypeName() string

	//
	Body() string

	//
	Marshal(*entities.Visit) (string, string, error)

	//
	Unmarshal([]byte) (*entities.Visit, error)
}