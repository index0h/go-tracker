package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type TrackHandler struct {
	trackManager *components.TrackManager
	uuid         dao.UUIDProviderInterface
}

func NewTrackHandler(trackManager *components.TrackManager, uuid dao.UUIDProviderInterface) *TrackHandler {
	return &TrackHandler{trackManager: trackManager, uuid: uuid}
}

func (handler *TrackHandler) Track(sessionID, clientID string, fields map[string]string) (*generated.Track, error) {
	visit, flashes, err := handler.trackManager.Track(handler.uuid.ToBytes(sessionID), clientID, fields)
	if err != nil {
		return nil, err
	}

	return &generated.Track{
		Visit:   handler.visitToThrift(visit),
		Flashes: handler.listFlashToThrift(flashes),
	}, nil
}

func (handler *TrackHandler) visitToThrift(input *entities.Visit) *generated.Visit {
	if input == nil {
		return nil
	}

	return &generated.Visit{
		VisitID:   handler.uuid.ToString(input.VisitID()),
		SessionID: handler.uuid.ToString(input.SessionID()),
		ClientID:  input.ClientID(),
		Timestamp: input.Timestamp(),
		Fields:    input.Fields(),
	}
}

func (handler *TrackHandler) listFlashToThrift(input []*entities.Flash) []*generated.Flash {
	if input == nil {
		return nil
	}

	result := make([]*generated.Flash, len(input))

	for i, value := range input {
		result[i] = &generated.Flash{
			FlashID:     handler.uuid.ToString(value.FlashID()),
			VisitID:     handler.uuid.ToString(value.VisitID()),
			EventID:     handler.uuid.ToString(value.EventID()),
			Timestamp:   value.Timestamp(),
			VisitFields: value.VisitFields(),
			EventFields: value.EventFields(),
		}
	}

	return result
}
