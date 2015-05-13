package uuid

type UuidMaker interface {
	Generate() Uuid
	ToString(Uuid) string
	ToBytes(string) Uuid
}
