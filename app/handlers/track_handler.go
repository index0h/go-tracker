package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	flashEntity "github.com/index0h/go-tracker/modules/flash/entity"
	"github.com/index0h/go-tracker/modules/track"
	visitEntity "github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share"
)

type TrackHandler struct {
	trackManager *track.Manager
	uuid         share.UUIDProviderInterface
}

func NewTrackHandler(trackManager *track.Manager, uuid share.UUIDProviderInterface) *TrackHandler {
	return &TrackHandler{trackManager: trackManager, uuid: uuid}
}

func (handler *TrackHandler) Track(sessionID, clientID string, fields map[string]string) (*generated.Track, error) {
	visit, flashes, err := handler.trackManager.Track(handler.uuid.FromString(sessionID), clientID, fields)
	if err != nil {
		return nil, err
	}

	return &generated.Track{
		Visit:   handler.visitToThrift(visit),
		Flashes: handler.listFlashToThrift(flashes),
	}, nil
}

func (handler *TrackHandler) visitToThrift(input *visitEntity.Visit) *generated.Visit {
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

func (handler *TrackHandler) listFlashToThrift(input []*flashEntity.Flash) []*generated.Flash {
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
