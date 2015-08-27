package share

type UUIDProviderInterface interface {
	// Generate new not empty UUID
	Generate() [16]byte

	// Converts UUID from bytes to string
	ToString([16]byte) string

	// Converts UUID from string to bytes
	ToBytes(string) [16]byte
}
