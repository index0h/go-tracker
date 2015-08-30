package uuid

import (
	"testing"

	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/share/types"
	"github.com/stretchr/testify/assert"
)

func Test_UUID_Interface(t *testing.T) {
	func(event share.UUIDProviderInterface) {}(New())
}

func Test_UUID_Generate_NotEmpty(t *testing.T) {
	var emptyUUID types.UUID

	checkUUID := new(UUID)

	assert.NotEqual(t, emptyUUID, checkUUID.Generate())
}

func Test_UUID_Generate_Duplicate(t *testing.T) {
	checkUUID := new(UUID)

	assert.NotEqual(t, checkUUID.Generate(), checkUUID.Generate())
}

func Test_UUID_ToStringEmpty(t *testing.T) {
	var emptyUUID types.UUID

	checkUUID := new(UUID)

	assert.Equal(t, "00000000-0000-0000-0000-000000000000", checkUUID.ToString(emptyUUID))
}

func Test_UUID_FromString_Empty(t *testing.T) {
	var emptyUUID types.UUID

	checkUUID := new(UUID)

	assert.Equal(t, emptyUUID, checkUUID.FromString("00000000-0000-0000-0000-000000000000"))
}

func Test_UUID_FromString_Panic(t *testing.T) {
	panicFunc := func() {
		New().FromString("Some wrong data here")
	}

	assert.Panics(t, panicFunc)
}

func Test_UUID_DoubleConvert(t *testing.T) {
	checkUUID := new(UUID)

	expectedBytes := checkUUID.Generate()
	expectedString := checkUUID.ToString(expectedBytes)

	actualBytes := checkUUID.FromString(expectedString)
	actualString := checkUUID.ToString(actualBytes)

	assert.Equal(t, expectedBytes, actualBytes)
	assert.Equal(t, expectedString, actualString)
}
