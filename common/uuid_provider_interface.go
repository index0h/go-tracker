package common

type UUIDProviderInterface interface {
	// Generate new not empty UUID
	Generate() UUID

	// Converts UUID from bytes to string
	ToString(UUID) string

	// Converts UUID from string to bytes
	ToBytes(string) UUID
}
