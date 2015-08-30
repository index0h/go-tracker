package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/modules/visit"
	"github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share"
)

type VisitHandler struct {
	visitManager *visit.Manager
	uuid         share.UUIDProviderInterface
}

func NewVisitHandler(visitManager *visit.Manager, uuid share.UUIDProviderInterface) *VisitHandler {
	return &VisitHandler{visitManager: visitManager, uuid: uuid}
}

func (handler *VisitHandler) FindVisitByID(visitID string) (*generated.Visit, error) {
	result, err := handler.visitManager.FindByID(handler.uuid.FromString(visitID))
	if err != nil {
		return nil, err
	}

	return handler.visitToThrift(result), nil
}

func (handler *VisitHandler) FindVisitAll(limit int64, offset int64) ([]*generated.Visit, error) {
	result, err := handler.visitManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil
}

func (handler *VisitHandler) FindVisitAllBySessionID(
	sessionID string,
	limit int64,
	offset int64,
) ([]*generated.Visit, error) {
	result, err := handler.visitManager.FindAllBySessionID(handler.uuid.FromString(sessionID), limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil
}

func (handler *VisitHandler) FindVisitAllByClientID(clientID string, limit, offset int64) ([]*generated.Visit, error) {
	result, err := handler.visitManager.FindAllByClientID(clientID, limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil
}

func (handler *VisitHandler) visitToThrift(input *entity.Visit) *generated.Visit {
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

func (handler *VisitHandler) listVisitToThrift(input []*entity.Visit) []*generated.Visit {
	if input == nil {
		return nil
	}

	result := make([]*generated.Visit, len(input))

	for i, value := range input {
		result[i] = handler.visitToThrift(value)
	}

	return result
}
