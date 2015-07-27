package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type VisitHandler struct {
	visitManager *components.VisitManager
	uuid         dao.UUIDProviderInterface
}

func NewVisitHandler(visitManager *components.VisitManager, uuid dao.UUIDProviderInterface) {
	return &VisitHandler{visitManager: visitManager, uuid: uuid}
}

func (handler *VisitHandler) FindByID(visitID string) (*tracker.Visit, error) {
	result, err := handler.visitManager.FindByID(handler.uuid.ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return handler.visitToThrift(result), nil
}

func (handler *VisitHandler) FindAll(limit int64, offset int64) ([]*tracker.Visit, error) {
	result, err := handler.visitManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil
}

func (handler *VisitHandler) FindAllBySessionID(sessionID string, limit, offset int64) ([]*tracker.Visit, error) {
	result, err := handler.visitManager.FindAllBySessionID(handler.uuid.ToBytes(sessionID), limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil
}

func (handler *VisitHandler) visitToThrift(input *entities.Visit) *tracker.Visit {
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

func (handler *VisitHandler) listVisitToThrift(input []*entities.Visit) []*tracker.Visit {
	if input == nil {
		return nil
	}

	result := make([]*tracker.Visit, len(input))

	for i, value := range input {
		result[i] = handler.visitToThrift(value)
	}

	return result
}
