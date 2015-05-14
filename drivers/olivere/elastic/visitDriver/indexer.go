package visitDriver

import "github.com/index0h/go-tracker/visit/entity"

type Indexer interface {
	Name(timestamp int64) string
	TypeName() string
	Body() string
	Marshal(*entity.Visit) (string, []byte, error)
	Unmarshal([]byte) (*entity.Visit, error)
}