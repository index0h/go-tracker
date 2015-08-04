package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Sirupsen/logrus"
	"github.com/index0h/go-servicelocator"
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/app/handlers"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao/dummy"
	"github.com/index0h/go-tracker/dao/elastic"
	"github.com/index0h/go-tracker/dao/memory"
	uuidDriver "github.com/index0h/go-tracker/dao/uuid"
	elasticClient "github.com/olivere/elastic"
	"github.com/index0h/go-tracker/dao"
)

var sl *servicelocator.ServiceLocator

func NewServiceLocator() *servicelocator.ServiceLocator {
	sl = servicelocator.New("tracker")

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	logger.Level = logrus.DebugLevel | logrus.InfoLevel | logrus.WarnLevel | logrus.ErrorLevel | logrus.FatalLevel

	sl.SetLogger(logger.WithField("service", "ServiceLocator"))

	sl.SetService("uuid", uuidDriver.New())
	sl.SetService("logger", logger)

	sl.SetService("visit_manager_logger", logger.WithField("service", "VisitManager"))
	sl.SetService("event_manager_logger", logger.WithField("service", "EventManager"))
	sl.SetService("flash_manager_logger", logger.WithField("service", "FlashManager"))
	sl.SetService("track_manager_logger", logger.WithField("service", "TrackManager"))
	sl.SetService("mark_manager_logger", logger.WithField("service", "MarkManager"))

	sl.SetConstructor("NewDummyVisitRepository", dummy.NewVisitRepository)
	sl.SetConstructor("NewDummyEventRepository", dummy.NewEventRepository)
	sl.SetConstructor("NewDummyFlashRepository", dummy.NewFlashRepository)
	sl.SetConstructor("NewDummyMarkRepository", dummy.NewMarkRepository)

	sl.SetConstructor("NewElasticClient", NewElasticClient)

	sl.SetConstructor("NewElasticVisitRepository", elastic.NewVisitRepository)
	sl.SetConstructor("NewElasticEventRepository", elastic.NewEventRepository)
	sl.SetConstructor("NewElasticFlashRepository", elastic.NewFlashRepository)

	sl.SetConstructor("NewMemoryVisitRepository", memory.NewVisitRepository)
	sl.SetConstructor("NewMemoryEventRepository", memory.NewEventRepository)

	sl.SetConstructor("NewVisitManager", components.NewVisitManager)
	sl.SetConstructor("NewEventManager", components.NewEventManager)
	sl.SetConstructor("NewFlashManager", components.NewFlashManager)
	sl.SetConstructor("NewMarkManager", components.NewMarkManager)
	sl.SetConstructor("NewTrackManager", components.NewTrackManager)

	sl.SetConfig(
		"visit_manager",
		"NewVisitManager",
		[]interface{}{"%visit_repository%", "%uuid%", "%visit_manager_logger%"},
	)
	sl.SetConfig(
		"event_manager",
		"NewEventManager",
		[]interface{}{"%event_repository%", "%uuid%", "%event_manager_logger%"},
	)
	sl.SetConfig(
		"flash_manager",
		"NewFlashManager",
		[]interface{}{"%flash_repository%", "%uuid%", "%flash_manager_logger%"},
	)
	sl.SetConfig(
		"mark_manager",
		"NewMarkManager",
		[]interface{}{"%mark_repository%", "%uuid%", "%mark_manager_logger%"},
	)
	sl.SetConfig(
		"track_manager",
		"NewTrackManager",
		[]interface{}{
			"%visit_manager%",
			"%event_manager%",
			"%flash_manager%",
			[]dao.ProcessorInterface{},
			"%uuid%",
			"%track_manager_logger%",
		},
	)

	sl.SetConstructor("NewVisitHandler", handlers.NewVisitHandler)
	sl.SetConstructor("NewEventHandler", handlers.NewEventHandler)
	sl.SetConstructor("NewFlashHandler", handlers.NewFlashHandler)
	sl.SetConstructor("NewTrackHandler", handlers.NewTrackHandler)
	sl.SetConstructor("NewMarkHandler", handlers.NewMarkHandler)

	sl.SetConfig("visit_handler", "NewVisitHandler", []interface{}{"%visit_manager%", "%uuid%"})
	sl.SetConfig("event_handler", "NewEventHandler", []interface{}{"%event_manager%", "%uuid%"})
	sl.SetConfig("flash_handler", "NewFlashHandler", []interface{}{"%flash_manager%", "%uuid%"})
	sl.SetConfig("track_handler", "NewTrackHandler", []interface{}{"%track_manager%", "%uuid%"})
	sl.SetConfig("mark_handler", "NewMarkHandler", []interface{}{"%mark_manager%", "%uuid%"})

	sl.SetConstructor("NewThriftService", NewThriftService)

	return sl
}

func NewElasticClient(host string, maxRetries int) (*elasticClient.Client, error) {
	return elasticClient.NewClient(elasticClient.SetURL(host), elasticClient.SetMaxRetries(maxRetries))
}

func NewThriftService(host string, bufferSize int) *thrift.TSimpleServer {
	processor := thrift.NewTMultiplexedProcessor()

	visitHandler, _ := sl.Get("visit_handler")
	eventHandler, _ := sl.Get("event_handler")
	flashHandler, _ := sl.Get("flash_handler")
	trackHandler, _ := sl.Get("track_handler")
	markHandler, _ := sl.Get("mark_handler")

	processor.RegisterProcessor("visit", generated.NewVisitServiceProcessor(visitHandler.(*handlers.VisitHandler)))
	processor.RegisterProcessor("event", generated.NewEventServiceProcessor(eventHandler.(*handlers.EventHandler)))
	processor.RegisterProcessor("flash", generated.NewFlashServiceProcessor(flashHandler.(*handlers.FlashHandler)))
	processor.RegisterProcessor("track", generated.NewTrackServiceProcessor(trackHandler.(*handlers.TrackHandler)))
	processor.RegisterProcessor("mark", generated.NewMarkServiceProcessor(markHandler.(*handlers.MarkHandler)))

	transport, err := thrift.NewTServerSocket(host)
	if err != nil {
		panic(err)
	}

	transportFactory := thrift.NewTBufferedTransportFactory(bufferSize)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	return thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
}
