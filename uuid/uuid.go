package uuid

type Uuid [16]byte

var emptyUuid Uuid

func IsUuidEmpty(u Uuid) bool {
	return (emptyUuid == u)
}

func NewEmpty() Uuid {
	return Uuid{}
}
