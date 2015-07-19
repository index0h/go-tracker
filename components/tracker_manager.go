package components

import (
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/entities"
	"log"
	"sort"
	"time"
	"fmt"
)

type TrackerManager struct {
	visitRepository dao.VisitRepositoryInterface
	eventRepository dao.EventRepositoryInterface
	flashRepository dao.FlashRepositoryInterface
	processors      []dao.ProcessorInterface
	uuid            dao.UUIDProviderInterface
	logger          *log.Logger
}

func NewTrackerManager(
	visitRepository dao.VisitRepositoryInterface,
	eventRepository dao.EventRepositoryInterface,
	flashRepository dao.FlashRepositoryInterface,
	processors []dao.ProcessorInterface,
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *TrackerManager {
	sort.Sort(ProcessorSorter{Data: processors})

	return &TrackerManager{
		visitRepository: visitRepository,
		eventRepository: eventRepository,
		flashRepository: flashRepository,
		processors:      processors,
		uuid:            uuid,
		logger:          logger,
	}
}

func (manager *TrackerManager) Track(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, flashes []*entities.Flash, err error) {
	visit, err = manager.createVisit(sessionID, clientID, fields)
	if err != nil {
		return visit, flashes, err
	}
	fmt.Printf("visit created:\n%+v", visit)

	if err := manager.visitRepository.Insert(visit); err != nil {
		return visit, flashes, err
	}
	fmt.Printf("visit inserted:\n%+v", visit)

	events, err := manager.eventRepository.FindAllByVisit(visit)
	if err != nil {
		return visit, flashes, err
	}
	fmt.Printf("events found:\n%+v", events)

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
			manager.flashRepository.Insert(flash)
			flashes = append(flashes, flash)
		}
	}
	fmt.Printf("result created:\n%+v", flashes)

	return visit, flashes, err
}

func (manager *TrackerManager) FindVisitByID(visitID [16]byte) (*entities.Visit, error) {
	return manager.visitRepository.FindByID(visitID)
}

func (manager *TrackerManager) FindVisitAll(limit int64, offset int64) ([]*entities.Visit, error) {
	return manager.visitRepository.FindAll(limit, offset)
}

func (manager *TrackerManager) FindVisitAllBySessionID(
	sessionID [16]byte,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	return manager.visitRepository.FindAllBySessionID(sessionID, limit, offset)
}

func (manager *TrackerManager) FindVisitAllByClientID(
	clientID string,
	limit int64,
	offset int64,
) ([]*entities.Visit, error) {
	return manager.visitRepository.FindAllByClientID(clientID, limit, offset)
}

func (manager *TrackerManager) InsertVisit(visit *entities.Visit) error {
	return manager.visitRepository.Insert(visit)
}

func (manager *TrackerManager) FindEventByID(eventID [16]byte) (*entities.Event, error) {
	return manager.eventRepository.FindByID(eventID)
}

func (manager *TrackerManager) FindEventAll(limit int64, offset int64) ([]*entities.Event, error) {
	return manager.eventRepository.FindAll(limit, offset)
}

func (manager *TrackerManager) InsertEvent(enabled bool, fields, filters entities.Hash) (*entities.Event, error) {
	event, err := entities.NewEvent(manager.uuid.Generate(), enabled, fields, filters)
	if err != nil {
		return nil, err
	}

	return event, manager.eventRepository.Insert(event)
}

func (manager *TrackerManager) UpdateEvent(event *entities.Event) (*entities.Event, error) {
	return event, manager.eventRepository.Update(event)
}

func (manager *TrackerManager) FindFlashByID(flashID [16]byte) (*entities.Flash, error) {
	return manager.flashRepository.FindByID(flashID)
}

func (manager *TrackerManager) FindFlashAll(limit int64, offset int64) ([]*entities.Flash, error) {
	return manager.flashRepository.FindAll(limit, offset)
}

func (manager *TrackerManager) FindFlashAllByVisitID(visitID [16]byte) ([]*entities.Flash, error) {
	return manager.flashRepository.FindAllByVisitID(visitID)
}

func (manager *TrackerManager) FindFlashAllByEventID(
	eventID [16]byte,
	limit int64,
	offset int64,
) ([]*entities.Flash, error) {
	return manager.flashRepository.FindAllByEventID(eventID, limit, offset)
}

func (manager *TrackerManager) createVisit(
	sessionID [16]byte,
	clientID string,
	fields entities.Hash,
) (visit *entities.Visit, err error) {
	if sessionID == [16]byte{} {
		sessionID = manager.uuid.Generate()
	} else {
		ok, err := manager.visitRepository.Verify(sessionID, clientID)

		if err != nil {
			return nil, err
		}

		if !ok {
			sessionID = manager.uuid.Generate()
			fields["warning:VisitManager"] = err.Error()
		}
	}

	return entities.NewVisit(manager.uuid.Generate(), time.Now().Unix(), sessionID, clientID, fields)
}
