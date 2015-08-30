package share

import "github.com/index0h/go-tracker/share/types"

type UUIDProviderInterface interface {
	// Generate new not empty UUID
	Generate() types.UUID

	// Converts UUID from bytes to string
	ToString(types.UUID) string

	// Converts UUID from string to bytes
	FromString(string) types.UUID
}
