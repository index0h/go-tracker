package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/index0h/go-servicelocator"
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/app/handlers"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/dummy"
	"github.com/index0h/go-tracker/dao/elastic"
	"github.com/index0h/go-tracker/dao/memory"
	uuidDriver "github.com/index0h/go-tracker/dao/uuid"
	elasticClient "github.com/olivere/elastic"
	"log"
	"os"
)

func NewServiceLocator() *servicelocator.ServiceLocator {
	result := servicelocator.New("tracker", "yaml")
	result.AddConfigPath("$HOME/.tracker")
	result.AddConfigPath("/etc/tracker")

	result.Set("NewLogger", NewLogger)
	result.Set("NewUUID", uuidDriver.New)
	result.Set("NewElasticClient", NewElasticClient)
	result.Set("NewDummyVisitRepository", dummy.NewVisitRepository())
	result.Set("NewElasticVisitRepository", elastic.NewVisitRepository)
	result.Set("NewMemoryVisitRepository", memory.NewVisitRepository)
	result.Set("NewVisitManager", components.NewVisitManager)
	result.Set("NewDummyEventRepository", dummy.NewEventRepository())
	result.Set("NewElasticEventRepository", elastic.NewEventRepository)
	result.Set("NewMemoryEventRepository", memory.NewEventRepository)
	result.Set("NewEventManager", components.NewEventManager)
	result.Set("NewDummyFlashRepository", dummy.NewFlashRepository())
	result.Set("NewElasticFlashRepository", elastic.NewFlashRepository)
	result.Set("NewFlashManager", components.NewFlashManager)
	result.Set("NewDummyMarkRepository", dummy.NewMarkRepository)
	result.Set("NewMarkManager", components.NewMarkManager)
	result.Set("NewTrackManager", NewTrackManager)
	result.Set("NewVisitHandler", handlers.NewVisitHandler)
	result.Set("NewEventHandler", handlers.NewEventHandler)
	result.Set("NewFlashHandler", handlers.NewFlashHandler)
	result.Set("NewTrackHandler", handlers.NewTrackHandler)
	result.Set("NewMarkHandler", handlers.NewMarkHandler)
	result.Set("NewThriftService", NewThriftService)

	return result
}

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "logger: ", log.Lshortfile)
}

func NewElasticClient(host string, maxRetries int) (*elasticClient.Client, error) {
	return elasticClient.NewClient(elasticClient.SetURL(host), elasticClient.SetMaxRetries(maxRetries))
}

func NewTrackManager(
	visitRepository dao.VisitRepositoryInterface,
	eventRepository dao.EventRepositoryInterface,
	flashRepository dao.FlashRepositoryInterface,
	processors string, // TODO
	uuid dao.UUIDProviderInterface,
	logger *log.Logger,
) *components.TrackManager {
	return components.NewTrackManager(visitRepository, eventRepository, flashRepository, nil, uuid, logger)
}

func NewThriftService(
	visitHandler *handlers.VisitHandler,
	eventHandler *handlers.EventHandler,
	flashHandler *handlers.FlashHandler,
	trackHandler *handlers.TrackHandler,
	markHandler *handlers.MarkHandler,
	host string,
	bufferSize int,
) *thrift.TSimpleServer {
	processor := thrift.NewTMultiplexedProcessor()
	processor.RegisterProcessor("visit", generated.NewVisitServiceProcessor(visitHandler))
	processor.RegisterProcessor("event", generated.NewEventServiceProcessor(eventHandler))
	processor.RegisterProcessor("flash", generated.NewFlashServiceProcessor(flashHandler))
	processor.RegisterProcessor("track", generated.NewTrackServiceProcessor(trackHandler))
	processor.RegisterProcessor("mark", generated.NewMarkServiceProcessor(markHandler))

	transport, err := thrift.NewTServerSocket(host)
	if err != nil {
		panic(err)
	}

	transportFactory := thrift.NewTBufferedTransportFactory(bufferSize)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	return thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
}
