package track

import (
	"sort"
	"time"

	eventPacket "github.com/index0h/go-tracker/modules/event"
	flashPacket "github.com/index0h/go-tracker/modules/flash"
	flashEntity "github.com/index0h/go-tracker/modules/flash/entity"
	visitPacket "github.com/index0h/go-tracker/modules/visit"
	visitEntity "github.com/index0h/go-tracker/modules/visit/entity"
	"github.com/index0h/go-tracker/share"
	"github.com/index0h/go-tracker/share/types"
)

type Manager struct {
	visitManager *visitPacket.Manager
	eventManager *eventPacket.Manager
	flashManager *flashPacket.Manager
	processors   []ProcessorInterface
	uuid         share.UUIDProviderInterface
	logger       share.LoggerInterface
}

func NewManager(
	visitManager *visitPacket.Manager,
	eventManager *eventPacket.Manager,
	flashManager *flashPacket.Manager,
	processors []ProcessorInterface,
	uuid share.UUIDProviderInterface,
	logger share.LoggerInterface,
) *Manager {
	sort.Sort(ProcessorSorter{Data: processors})

	return &Manager{
		visitManager: visitManager,
		eventManager: eventManager,
		flashManager: flashManager,
		processors:   processors,
		uuid:         uuid,
		logger:       logger,
	}
}

func (manager *Manager) Track(
	sessionID types.UUID,
	clientID string,
	fields types.Hash,
) (visit *visitEntity.Visit, flashes []*flashEntity.Flash, err error) {
	visit, err = manager.visitManager.CreateVisit(sessionID, clientID, fields)
	if err != nil {
		return visit, flashes, err
	}

	if err := manager.visitManager.InsertVisit(visit); err != nil {
		return visit, flashes, err
	}

	events, err := manager.eventManager.FindAllByFields(visit.Fields())
	if err != nil {
		return visit, flashes, err
	}

	timestamp := time.Now().Unix()

	for _, event := range events {
		flash, err := flashEntity.NewFlash(
			manager.uuid.Generate(),
			visit.VisitID(),
			event.EventID(),
			timestamp,
			visit.Fields(),
			event.Fields(),
		)
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
