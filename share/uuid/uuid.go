package uuid

import (
	"github.com/index0h/go-tracker/share/types"
	satoriUUID "github.com/satori/go.uuid"
)

type UUID struct{}

// Check that uuid is empty
func New() *UUID {
	return &UUID{}
}

// Generate new not empty uuid
func (uuid *UUID) Generate() types.UUID {
	var result types.UUID

	copy(result[:], satoriUUID.NewV4().Bytes())

	return result
}

// Converts uuid from bytes to string
func (uuid *UUID) ToString(uuidBytes types.UUID) string {
	// Force set 16 bytes, so there is no errors could be in satori UUID
	result, _ := satoriUUID.FromBytes(uuidBytes[:])

	return result.String()
}

// Converts uuid from string to bytes
func (uuid *UUID) FromString(uuidString string) types.UUID {
	if uuidString == "" {
		return types.UUID{}
	}

	uuidResult, err := satoriUUID.FromString(uuidString)

	if err != nil {
		panic("Invalid UUID string")
	}

	var result types.UUID

	copy(result[:], uuidResult.Bytes())

	return result
}
