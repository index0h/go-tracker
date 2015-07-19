package main

import (
	"errors"
	"github.com/index0h/go-tracker/app/tracker"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/entities"
)

type ThriftHandler struct {
	trackerManager *components.TrackerManager
}

func (handler *ThriftHandler) Track(sessionID, clientID string, fields map[string]string) (*tracker.TrackResponse, error) {
	visit, flashes, err := handler.trackerManager.Track(getUUID().ToBytes(sessionID), clientID, fields)
	if err != nil {
		return nil, err
	}

	return &tracker.TrackResponse{
		Visit:   entitiesVisitToTrackerVisit(visit),
		Flashes: listEntitiesFlashToTrackerFlash(flashes),
	}, nil
}

func (handler *ThriftHandler) FindVisitByID(visitID string) (*tracker.Visit, error) {
	result, err := handler.trackerManager.FindVisitByID(getUUID().ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return entitiesVisitToTrackerVisit(result), nil
}

func (handler *ThriftHandler) FindVisitAll(limit int64, offset int64) ([]*tracker.Visit, error) {
	result, err := handler.trackerManager.FindVisitAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesVisitToTrackerVisit(result), nil
}

func (handler *ThriftHandler) FindVisitAllBySessionID(sessionID string, limit, offset int64) ([]*tracker.Visit, error) {
	result, err := handler.trackerManager.FindVisitAllBySessionID(getUUID().ToBytes(sessionID), limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesVisitToTrackerVisit(result), nil
}

func (handler *ThriftHandler) FindVisitAllByClientID(clientID string, limit, offset int64) ([]*tracker.Visit, error) {
	result, err := handler.trackerManager.FindVisitAllByClientID(clientID, limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesVisitToTrackerVisit(result), nil
}

func (handler *ThriftHandler) FindEventByID(eventID string) (*tracker.Event, error) {
	result, err := handler.trackerManager.FindEventByID(getUUID().ToBytes(eventID))
	if err != nil {
		return nil, err
	}

	return entitiesEventToTrackerEvent(result), nil
}

func (handler *ThriftHandler) FindEventAll(limit int64, offset int64) ([]*tracker.Event, error) {
	result, err := handler.trackerManager.FindEventAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesEventToTrackerEvent(result), nil
}

func (handler *ThriftHandler) InsertEvent(enabled bool, fields, filters map[string]string) (*tracker.Event, error) {
	event, err := handler.trackerManager.InsertEvent(enabled, fields, filters)

	if err != nil {
		return nil, err
	}

	return entitiesEventToTrackerEvent(event), err
}

func (handler *ThriftHandler) UpdateEvent(event *tracker.Event) (*tracker.Event, error) {
	eventModel, err := trackerEventToEntitiesEvent(event)

	if err != nil {
		return event, err
	}

	eventModel, err = handler.trackerManager.UpdateEvent(eventModel)

	if err != nil {
		return event, err
	}

	return entitiesEventToTrackerEvent(eventModel), err
}

func (handler *ThriftHandler) FindFlashByID(flashID string) (*tracker.Flash, error) {
	result, err := handler.trackerManager.FindFlashByID(getUUID().ToBytes(flashID))
	if err != nil {
		return nil, err
	}

	return entitiesFlashToTrackerFlash(result), nil
}

func (handler *ThriftHandler) FindFlashAll(limit int64, offset int64) ([]*tracker.Flash, error) {
	result, err := handler.trackerManager.FindFlashAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesFlashToTrackerFlash(result), nil
}

func (handler *ThriftHandler) FindFlashAllByVisitID(visitID string) ([]*tracker.Flash, error) {
	result, err := handler.trackerManager.FindFlashAllByVisitID(getUUID().ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return listEntitiesFlashToTrackerFlash(result), nil
}

func (handler *ThriftHandler) FindFlashAllByEventID(eventID string, limit, offset int64) ([]*tracker.Flash, error) {
	result, err := handler.trackerManager.FindFlashAllByEventID(getUUID().ToBytes(eventID), limit, offset)
	if err != nil {
		return nil, err
	}

	return listEntitiesFlashToTrackerFlash(result), nil
}

func entitiesVisitToTrackerVisit(input *entities.Visit) *tracker.Visit {
	if input == nil {
		return nil
	}

	return &tracker.Visit{
		VisitID:   getUUID().ToString(input.VisitID()),
		SessionID: getUUID().ToString(input.SessionID()),
		ClientID:  input.ClientID(),
		Timestamp: input.Timestamp(),
		Fields:    input.Fields(),
	}
}

func listEntitiesVisitToTrackerVisit(input []*entities.Visit) []*tracker.Visit {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Visit, len(input))

	for i, value := range input {
		result[i] = entitiesVisitToTrackerVisit(value)
	}

	return result
}

func entitiesEventToTrackerEvent(input *entities.Event) *tracker.Event {
	if input == nil {
		return nil
	}

	return &tracker.Event{
		EventID: getUUID().ToString(input.EventID()),
		Enabled: input.Enabled(),
		Fields:  input.Fields(),
		Filters: input.Filters(),
	}
}

func listEntitiesEventToTrackerEvent(input []*entities.Event) []*tracker.Event {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Event, len(input))

	for i, value := range input {
		result[i] = entitiesEventToTrackerEvent(value)
	}

	return result
}

func entitiesFlashToTrackerFlash(input *entities.Flash) *tracker.Flash {
	if input == nil {
		return nil
	}

	return &tracker.Flash{
		FlashID:     getUUID().ToString(input.FlashID()),
		VisitID:     getUUID().ToString(input.VisitID()),
		EventID:     getUUID().ToString(input.EventID()),
		Timestamp:   input.Timestamp(),
		VisitFields: input.VisitFields(),
		EventFields: input.EventFields(),
	}
}

func listEntitiesFlashToTrackerFlash(input []*entities.Flash) []*tracker.Flash {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Flash, len(input))

	for i, value := range input {
		result[i] = entitiesFlashToTrackerFlash(value)
	}

	return result
}

func trackerEventToEntitiesEvent(input *tracker.Event) (*entities.Event, error) {
	if input == nil {
		return nil, errors.New("input event must not be nil")
	}

	return entities.NewEvent(getUUID().ToBytes(input.EventID), input.Enabled, input.Fields, input.Filters)
}
