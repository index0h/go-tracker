package types

type UUID [16]byte

func (uuid *UUID) IsEmpty() bool {
	return uuid == [16]byte{}
}
