package elastic

type InternalRepositoryInterface interface {
	Refresh(scenario int) bool

	GetIndexName(scenario int) string
	GetTypeName() string

	GetEntityID(interface{}) (string, error)
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte) (interface{}, error)
}
