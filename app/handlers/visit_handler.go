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

func NewVisitHandler(visitManager *components.VisitManager, uuid dao.UUIDProviderInterface) *VisitHandler {
	return &VisitHandler{visitManager: visitManager, uuid: uuid}
}

func (handler *VisitHandler) FindVisitByID(visitID string) (*generated.Visit, error) {
	/*result, err := handler.visitManager.FindByID(handler.uuid.ToBytes(visitID))
	if err != nil {
		return nil, err
	}

	return handler.visitToThrift(result), nil*/

	panic("NOT IMPLEMENTED")
}

func (handler *VisitHandler) FindVisitAll(limit int64, offset int64) ([]*generated.Visit, error) {
	/*result, err := handler.visitManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil*/

	panic("NOT IMPLEMENTED")
}

func (handler *VisitHandler) FindVisitAllBySessionID(sessionID string, limit, offset int64) ([]*generated.Visit, error) {
	/*result, err := handler.visitManager.FindAllBySessionID(handler.uuid.ToBytes(sessionID), limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil*/

	panic("NOT IMPLEMENTED")
}

func (handler *VisitHandler) FindVisitAllByClientID(clientID string, limit, offset int64) ([]*generated.Visit, error) {
	/*result, err := handler.visitManager.FindAllByClientID(clientID, limit, offset)
	if err != nil {
		return nil, err
	}

	return handler.listVisitToThrift(result), nil*/

	panic("NOT IMPLEMENTED")
}

func (handler *VisitHandler) visitToThrift(input *entities.Visit) *generated.Visit {
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

func (handler *VisitHandler) listVisitToThrift(input []*entities.Visit) []*generated.Visit {
	if input == nil {
		return nil
	}

	result := make([]*generated.Visit, len(input))

	for i, value := range input {
		result[i] = handler.visitToThrift(value)
	}

	return result
}
