package common

type UUID [16]byte

func IsUUIDEmpty(uuid UUID) bool {
	return (uuid == UUID{})
}
