package entities

import "errors"

type hashListFields map[string][]string

func (hash *hashListFields) Copy() hashListFields {
	result := make(Hash, len(*hash))

	for key, value := range *hash {
		result[key] = make([]string, len(value))
		copy(result[key], value)
	}

	return result
}

type Mark struct {
	markID     [16]byte
	clientID   string
	sessionID  [16]byte
	fields     Hash
	listFields hashListFields
}

func NewMark(
	markID [16]byte,
	clientID string,
	sessionID [16]byte,
	fields Hash,
	listFields hashListFields,
) (*Mark, error) {
	if markID == [16]byte{} {
		return nil, errors.New("Empty markID is not allowed")
	}

	if (clientID == "") && (sessionID == [16]byte{}) {
		return nil, errors.New("Empty clientID and sessionID not allowed")
	}

	return &Mark{
		markID:     markID,
		clientID:   clientID,
		sessionID:  sessionID,
		fields:     fields.Copy(),
		listFields: listFields.Copy(),
	}, nil
}

func (mark *Mark) MarkID() [16]byte {
	return mark.markID
}

func (mark *Mark) ClientID() string {
	return mark.clientID
}

func (mark *Mark) SessionID() [16]byte {
	return mark.sessionID
}

func (mark *Mark) VisitFields() Hash {
	return mark.fields.Copy()
}

func (mark *Mark) ListFields() hashListFields {
	return mark.listFields.Copy()
}
