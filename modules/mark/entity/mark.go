package entity

import (
	"errors"
	"github.com/index0h/go-tracker/share/types"
)

type hashListFields map[string][]string

func (hash *hashListFields) Copy() hashListFields {
	result := make(hashListFields, len(*hash))

	for key, value := range *hash {
		result[key] = make([]string, len(value))
		copy(result[key], value)
	}

	return result
}

type Mark struct {
	markID     types.UUID
	sessionID  types.UUID
	clientID   string
	fields     types.Hash
	listFields hashListFields
}

func NewMark(
	markID types.UUID,
	sessionID types.UUID,
	clientID string,
	fields types.Hash,
	listFields hashListFields,
) (*Mark, error) {
	if markID.IsEmpty() {
		return nil, errors.New("Empty markID is not allowed")
	}

	if (clientID == "") && sessionID.IsEmpty() {
		return nil, errors.New("Empty clientID and sessionID not allowed")
	}

	return &Mark{
		markID:     markID,
		sessionID:  sessionID,
		clientID:   clientID,
		fields:     fields.Copy(),
		listFields: listFields.Copy(),
	}, nil
}

func (mark *Mark) MarkID() types.UUID {
	return mark.markID
}

func (mark *Mark) ClientID() string {
	return mark.clientID
}

func (mark *Mark) SessionID() types.UUID {
	return mark.sessionID
}

func (mark *Mark) VisitFields() types.Hash {
	return mark.fields.Copy()
}

func (mark *Mark) ListFields() hashListFields {
	return mark.listFields.Copy()
}
