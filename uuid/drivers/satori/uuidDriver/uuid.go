package uuidDriver

import (
	"github.com/index0h/go-tracker/uuid"
	satoriUuid "github.com/satori/go.uuid"
)

type Uuid struct {
}

func (u Uuid) Generate() uuid.Uuid {
	var result uuid.Uuid
	copy(result[:], satoriUuid.NewV4().Bytes())
	return result
}

func (u Uuid) ToString(uuidBytes uuid.Uuid) string {
	result, err := satoriUuid.FromBytes(uuidBytes[:])
	if err != nil {
		panic("Invalid UUID bytes")
	}
	return result.String()
}

func (u Uuid) ToBytes(uuidString string) uuid.Uuid {
	var result uuid.Uuid
	uuidResult, err := satoriUuid.FromString(uuidString)
	if err != nil {
		panic("Invalid UUID string")
	}
	copy(result[:], uuidResult.Bytes())
	return result
}
