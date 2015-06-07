package uuid

import (
	"github.com/index0h/go-tracker/common"
	satoriUUID "github.com/satori/go.uuid"
)

type UUID struct {
}

// Check that uuid is empty
func New() *UUID {
	return &UUID{}
}

// Generate new not empty uuid
func (uuid *UUID) Generate() common.UUID {
	var result common.UUID

	copy(result[:], satoriUUID.NewV4().Bytes())

	return result
}

// Converts uuid from bytes to string
func (uuid *UUID) ToString(uuidBytes common.UUID) string {
	result, err := satoriUUID.FromBytes(uuidBytes[:])

	if err != nil {
		panic("Invalid UUID bytes")
	}

	return result.String()
}

// Converts uuid from string to bytes
func (uuid *UUID) ToBytes(uuidString string) common.UUID {
	uuidResult, err := satoriUUID.FromString(uuidString)

	if err != nil {
		panic("Invalid UUID string")
	}

	var result common.UUID

	copy(result[:], uuidResult.Bytes())

	return result
}
