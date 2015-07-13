package main

import (
	"errors"
	"github.com/index0h/go-tracker/app/tracker"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"log"
)

type ThriftHandler struct {
	visitManager   *components.VisitManager
	eventManager   *components.EventManager
	flashManager   *components.FlashManager
	trackerManager *components.TrackerManager
	uuid           dao.UUIDProviderInterface
	logger         *log.Logger
}

func (handler *ThriftHandler) Track(sessionID, clientID string, fields map[string]string) ([]*tracker.Flash, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindVisitByID(visitID string) (*tracker.Visit, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindVisitAll(limit int64, offset int64) ([]*tracker.Visit, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindVisitAllBySessionID(sessionID string, limit, offset int64) ([]*tracker.Visit, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindVisitAllByClientID(clientID string, limit, offset int64) ([]*tracker.Visit, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindEventByID(eventID string) (*tracker.Event, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindEventAll(limit int64, offset int64) ([]*tracker.Event, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) InsertEvent(event *tracker.Event) error {
	return nil
}

func (handler *ThriftHandler) UpdateEvent(event *tracker.Event) error {
	return nil
}

func (handler *ThriftHandler) FindFlashByID(flashID string) (*tracker.Flash, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindFlashAll(limit int64, offset int64) ([]*tracker.Flash, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindFlashAllByVisitID(visitID string, limit, offset int64) ([]*tracker.Flash, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (handler *ThriftHandler) FindFlashAllByEventID(eventID string, limit, offset int64) ([]*tracker.Flash, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}
