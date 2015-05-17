package uuid

type Maker interface {
	// Generate new not empty uuid
	Generate() UUID

	// Converts uuid from bytes to string
	ToString(UUID) string

	// Converts uuid from string to bytes
	ToBytes(string) UUID
}
