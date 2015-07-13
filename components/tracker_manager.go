package components

import (
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	"log"
	"sort"
	"time"
)

type TrackerManager struct {
	visitManager *VisitManager
	eventManager *EventManager
	flashManager *FlashManager
	processors   []dao.ProcessorInterface
	uuid         dao.UUIDProviderInterface
	logger       *log.Logger
}

func NewTrackerManager(
	visitManager *VisitManager,
	eventManager *EventManager,
	flashManager *FlashManager,
	processors []dao.ProcessorInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *TrackerManager {
	sort.Sort(ProcessorSorter{Data: processors})

	return &TrackerManager{
		visitManager: visitManager,
		eventManager: eventManager,
		flashManager: flashManager,
		processors:   processors,
		uuid:         uuid,
		logger:       logger,
	}
}

func (manager *TrackerManager) Track(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (result []*entities.Flash, err error) {
	visit, err := manager.visitManager.Insert(sessionID, clientID, fields)
	if err != nil {
		return result, err
	}

	events, err := manager.eventManager.FindAllByVisit(visit)
	if err != nil {
		return result, err
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
			result = append(result, flash)
		}
	}

	return result, err
}
