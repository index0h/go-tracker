package components

import (
	"sort"
	"time"

	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
)

type TrackManager struct {
	visitManager *VisitManager
	eventManager *EventManager
	flashManager *FlashManager
	processors   []dao.ProcessorInterface
	uuid         dao.UUIDProviderInterface
	logger       dao.LoggerInterface
}

func NewTrackManager(
	visitManager *VisitManager,
	eventManager *EventManager,
	flashManager *FlashManager,
	processors []dao.ProcessorInterface,
	uuid dao.UUIDProviderInterface,
	logger dao.LoggerInterface,
) *TrackManager {
	sort.Sort(ProcessorSorter{Data: processors})

	return &TrackManager{
		visitManager: visitManager,
		eventManager: eventManager,
		flashManager: flashManager,
		processors:   processors,
		uuid:         uuid,
		logger:       logger,
	}
}

func (manager *TrackManager) Track(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, flashes []*entities.Flash, err error) {
	visit, err = manager.visitManager.CreateVisit(sessionID, clientID, fields)
	if err != nil {
		return visit, flashes, err
	}

	if err := manager.visitManager.InsertVisit(visit); err != nil {
		return visit, flashes, err
	}

	events, err := manager.eventManager.FindAllByVisit(visit)
	if err != nil {
		return visit, flashes, err
	}

	for _, event := range events {
		flash, err := entities.NewFlash(manager.uuid.Generate(), time.Now().Unix(), visit, event)
		if err != nil {
			break
		}

		for _, processor := range manager.processors {
			flash = processor.Process(flash, event, visit)
			if flash == nil {
				break
			}
		}

		if flash != nil {
			manager.flashManager.Insert(flash)
			flashes = append(flashes, flash)
		}
	}

	return visit, flashes, err
}
