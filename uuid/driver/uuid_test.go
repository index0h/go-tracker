package driver

import (
	"testing"

	interfaceUUID "github.com/index0h/go-tracker/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	func(event interfaceUUID.Maker) {}(&UUID{})
}

func TestGenerateNotEmpty(t *testing.T) {
	var emptyUUID interfaceUUID.UUID

	checkUUID := new(UUID)

	assert.NotEqual(t, emptyUUID, checkUUID.Generate())
}

func TestGenerateDuplicate(t *testing.T) {
	checkUUID := new(UUID)

	assert.NotEqual(t, checkUUID.Generate(), checkUUID.Generate())
}

func TestToStringEmpty(t *testing.T) {
	var emptyUUID interfaceUUID.UUID

	checkUUID := new(UUID)

	assert.Equal(t, "00000000-0000-0000-0000-000000000000", checkUUID.ToString(emptyUUID))
}

func TestToBytesEmpty(t *testing.T) {
	var emptyUUID interfaceUUID.UUID

	checkUUID := new(UUID)

	assert.Equal(t, emptyUUID, checkUUID.ToBytes("00000000-0000-0000-0000-000000000000"))
}

func TestDoubleConvert(t *testing.T) {
	checkUUID := new(UUID)

	expectedBytes := checkUUID.Generate()
	expectedString := checkUUID.ToString(expectedBytes)

	actualBytes := checkUUID.ToBytes(expectedString)
	actualString := checkUUID.ToString(actualBytes)

	assert.Equal(t, expectedBytes, actualBytes)
	assert.Equal(t, expectedString, actualString)
}
