package entities

type Mark struct {
	markID     [16]byte
	clientID   [16]byte
	fields     Hash
	listFields map[string][]string
}
