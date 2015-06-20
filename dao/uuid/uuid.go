package uuid

import (
	satoriUUID "github.com/satori/go.uuid"
)

type UUID struct{}

// Check that uuid is empty
func New() *UUID {
	return &UUID{}
}

// Generate new not empty uuid
func (uuid *UUID) Generate() [16]byte {
	var result [16]byte

	copy(result[:], satoriUUID.NewV4().Bytes())

	return result
}

// Converts uuid from bytes to string
func (uuid *UUID) ToString(uuidBytes [16]byte) string {
	// Force set 16 bytes, so there is no errors could be in satori UUID
	result, _ := satoriUUID.FromBytes(uuidBytes[:])

	return result.String()
}

// Converts uuid from string to bytes
func (uuid *UUID) ToBytes(uuidString string) [16]byte {
	uuidResult, err := satoriUUID.FromString(uuidString)

	if err != nil {
		panic("Invalid UUID string")
	}

	var result [16]byte

	copy(result[:], uuidResult.Bytes())

	return result
}
