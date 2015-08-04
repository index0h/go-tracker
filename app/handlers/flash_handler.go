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

func NewFlashHandler(flashManager *components.FlashManager, uuid dao.UUIDProviderInterface) *FlashHandler {
	return &FlashHandler{flashManager: flashManager, uuid: uuid}
}

func (handler *FlashHandler) FindFlashByID(flashID string) (*generated.Flash, error) {
	result, err := handler.flashManager.FindByID(handler.uuid.ToBytes(flashID))
	if err != nil {
		return nil, err
	}

	return handler.flashToThrift(result), nil
}

func (handler *FlashHandler) FindFlashAll(limit int64, offset int64) ([]*generated.Flash, error) {
	result, err := handler.flashManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) FindFlashAllByVisitID(visitID string) ([]*generated.Flash, error) {
	result, err := handler.flashManager.FindAllByVisitID(handler.uuid.ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) FindFlashAllByEventID(eventID string, limit, offset int64) ([]*generated.Flash, error) {
	result, err := handler.flashManager.FindAllByEventID(handler.uuid.ToBytes(eventID), limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listFlashToThrift(result), nil
}

func (handler *FlashHandler) flashToThrift(input *entities.Flash) *generated.Flash {
	if input == nil {
		return nil
	}

	return &generated.Flash{
		FlashID:     handler.uuid.ToString(input.FlashID()),
		VisitID:     handler.uuid.ToString(input.VisitID()),
		EventID:     handler.uuid.ToString(input.EventID()),
		Timestamp:   input.Timestamp(),
		VisitFields: input.VisitFields(),
		EventFields: input.EventFields(),
	}
}

func (handler *FlashHandler) listFlashToThrift(input []*entities.Flash) []*generated.Flash {
	if input == nil {
		return nil
	}

	result := make([]*generated.Flash, len(input))

	for i, value := range input {
		result[i] = handler.flashToThrift(value)
	}

	return result
}
