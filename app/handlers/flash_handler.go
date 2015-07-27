package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type FlashHandler struct {
	flashManager *components.FlashManager
	uuid         dao.UUIDProviderInterface
}

func NewFlashHandler(flashManager *components.FlashManager, uuid dao.UUIDProviderInterface) {
	return &FlashHandler{flashManager: flashManager, uuid: uuid}
}

func (handler *FlashHandler) FindByID(flashID string) (*tracker.Flash, error) {
	result, err := handler.flashManager.FindByID(handler.uuid.ToBytes(flashID))
	if err != nil {
		return nil, err
	}

	return handler.flashToThrift(result), nil
}

func (handler *FlashHandler) FindAll(limit int64, offset int64) ([]*tracker.Flash, error) {
	result, err := handler.flashManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) FindAllByVisitID(visitID string) ([]*tracker.Flash, error) {
	result, err := handler.flashManager.FindAllByVisitID(handler.uuid.ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) FindAllByEventID(eventID string, limit, offset int64) ([]*tracker.Flash, error) {
	result, err := handler.flashManager.FindAllByEventID(handler.uuid.ToBytes(eventID), limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) flashToThrift(input *entities.Flash) *tracker.Flash {
	if input == nil {
		return nil
	}

	return &tracker.Flash{
		FlashID:     handler.uuid.ToString(input.FlashID()),
		VisitID:     handler.uuid.ToString(input.VisitID()),
		EventID:     handler.uuid.ToString(input.EventID()),
		Timestamp:   input.Timestamp(),
		VisitFields: input.VisitFields(),
		EventFields: input.EventFields(),
	}
}

func (handler *FlashHandler) listFlashToThrift(input []*entities.Flash) []*tracker.Flash {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Flash, len(input))

	for i, value := range input {
		result[i] = handler.flashToThrift(value)
	}

	return result
}
