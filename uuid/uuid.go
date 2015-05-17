package uuid

type UUID [16]byte

var emptyUUID UUID

// Check that uuid is empty
func IsUUIDEmpty(u UUID) bool {
	return (emptyUUID == u)
}

// Creates empty uuid
func NewEmpty() UUID {
	return UUID{}
}
