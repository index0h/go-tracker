package handlers

import (
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/modules/mark"
	"github.com/index0h/go-tracker/share"
)

type MarkHandler struct {
	markManager *mark.Manager
	uuid        share.UUIDProviderInterface
}

func NewMarkHandler(markManager *mark.Manager, uuid share.UUIDProviderInterface) *MarkHandler {
	return &MarkHandler{markManager: markManager, uuid: uuid}
}

func (handler *MarkHandler) FindMarkByID(markID string) (*generated.Mark, error) {
	/*_, err := handler.markManager.FindByID(handler.uuid.ToBytes(markID))
	if err != nil {
		return nil, err
	}*/

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) FindMarkByClientID(clientID string) (*generated.Mark, error) {
	/*_, err := handler.markManager.FindByClientID(clientID)
	if clientID != nil {
		return nil, err
	}*/

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) FindMarkAll(limit int64, offset int64) ([]*generated.Mark, error) {
	/*_, err := handler.markManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}*/

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) UpdateMark(mark *generated.Mark) (*generated.Mark, error) {
	panic("NOT IMPLEMENTED")
}
