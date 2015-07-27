package handlers

import (
	"github.com/index0h/go-tracker/app/thrift/tracker"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
)

type MarkHandler struct {
	markManager *components.MarkManager
	uuid        dao.UUIDProviderInterface
}

func NewMarkHandler(markManager *components.MarkManager, uuid dao.UUIDProviderInterface) {
	return &MarkHandler{markManager: markManager, uuid: uuid}
}

func (handler *MarkHandler) FindByID(markID string) (*tracker.Mark, error) {
	_, err := handler.markManager.FindByID(handler.uuid.ToBytes(markID))
	if err != nil {
		return nil, err
	}

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) FindByClientID(clientID string) (*tracker.Mark, error) {
	_, err := handler.markManager.FindByClientID(clientID)
	if clientID != nil {
		return nil, err
	}

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) FindAll(limit int64, offset int64) ([]*tracker.Mark, error) {
	_, err := handler.markManager.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	panic("NOT IMPLEMENTED")
}

func (handler *MarkHandler) Update(event *tracker.Mark) (*tracker.Mark, error) {
	panic("NOT IMPLEMENTED")
}
