package uuidDriver

import (
	interfaceUUID "github.com/index0h/go-tracker/uuid"
	satoriUUID "github.com/satori/go.uuid"
)

type UUID struct {
}

// Generate new not empty uuid
func (uuid *UUID) Generate() interfaceUUID.UUID {
	var result interfaceUUID.UUID

	copy(result[:], satoriUUID.NewV4().Bytes())

	return result
}

// Converts uuid from bytes to string
func (uuid *UUID) ToString(uuidBytes interfaceUUID.UUID) string {
	result, err := satoriUUID.FromBytes(uuidBytes[:])

	if err != nil {
		panic("Invalid UUID bytes")
	}

	return result.String()
}

// Converts uuid from string to bytes
func (uuid *UUID) ToBytes(uuidString string) interfaceUUID.UUID {
	uuidResult, err := satoriUUID.FromString(uuidString)

	if err != nil {
		panic("Invalid UUID string")
	}

	var result interfaceUUID.UUID

	copy(result[:], uuidResult.Bytes())

	return result
}
