package handlers

import (
	"errors"
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type EventHandler struct {
	eventManager *components.EventManager
	uuid         dao.UUIDProviderInterface
}

func NewEventHandler(eventManager *components.EventManager, uuid dao.UUIDProviderInterface) *EventHandler {
	return &EventHandler{eventManager: eventManager, uuid: uuid}
}

func (handler *EventHandler) FindEventByID(eventID string) (*generated.Event, error) {
	result, err := handler.eventManager.FindByID(handler.uuid.ToBytes(eventID))
	if err != nil {
		return nil, err
	}

	return handler.eventToThrift(result), nil
}

func (handler *EventHandler) FindEventAll(limit int64, offset int64) ([]*generated.Event, error) {
	result, err := handler.eventManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listEventToThrift(result), nil
}

func (handler *EventHandler) InsertEvent(enabled bool, fields, filters map[string]string) (*generated.Event, error) {
	event, err := handler.eventManager.Insert(enabled, fields, filters)
	if err != nil {
		return nil, err
	}

	return handler.eventToThrift(event), nil
}

func (handler *EventHandler) UpdateEvent(event *generated.Event) (*generated.Event, error) {
	eventModel, err := handler.thriftToEvent(event)
	if err != nil {
		return event, err
	}

	err = handler.eventManager.Update(eventModel)
	if err != nil {
		return event, err
	}

	return handler.eventToThrift(eventModel), err
}

func (handler *EventHandler) eventToThrift(input *entities.Event) *generated.Event {
	if input == nil {
		return nil
	}

	return &generated.Event{
		EventID: handler.uuid.ToString(input.EventID()),
		Enabled: input.Enabled(),
		Fields:  input.Fields(),
		Filters: input.Filters(),
	}
}

func (handler *EventHandler) listEventToThrift(input []*entities.Event) []*generated.Event {
	if input == nil {
		return nil
	}

	result := make([]*generated.Event, len(input))

	for i, value := range input {
		result[i] = handler.eventToThrift(value)
	}

	return result
}

func (handler *EventHandler) thriftToEvent(input *generated.Event) (*entities.Event, error) {
	if input == nil {
		return nil, errors.New("input event must not be nil")
	}

	return entities.NewEvent(handler.uuid.ToBytes(input.EventID), input.Enabled, input.Fields, input.Filters)
}
