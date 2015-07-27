package handlers

import (
	"errors"
	"github.com/index0h/go-tracker/app/thrift/tracker"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type TrackHandler struct {
	trackManager *components.TrackManager
	uuid         dao.UUIDProviderInterface
}

func NewTrackHandler(trackManager *components.TrackManager, uuid dao.UUIDProviderInterface) {
	return &TrackHandler{trackManager: trackManager, uuid: uuid}
}

func (handler *TrackHandler) Track(sessionID, clientID string, fields map[string]string) (*tracker.Track, error) {
	visit, flashes, err := handler.trackManager.Track(handler.uuid.ToBytes(sessionID), clientID, fields)
	if err != nil {
		return nil, err
	}

	return &tracker.Track{
		Visit:   handler.visitToThrift(visit),
		Flashes: handler.listFlashToThrift(flashes),
	}, nil
}

func (handler *TrackHandler) visitToThrift(input *entities.Visit) *tracker.Visit {
	if input == nil {
		return nil
	}

	return &tracker.Visit{
		VisitID:   handler.uuid.ToString(input.VisitID()),
		SessionID: handler.uuid.ToString(input.SessionID()),
		ClientID:  input.ClientID(),
		Timestamp: input.Timestamp(),
		Fields:    input.Fields(),
	}
}

func (handler *TrackHandler) listFlashToThrift(input []*entities.Flash) []*tracker.Flash {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Flash, len(input))

	for i, value := range input {
		result[i] = &tracker.Flash{
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
