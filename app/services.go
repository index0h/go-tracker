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
	"strconv"
)

func NewServiceLocator() *servicelocator.ServiceLocator {
	result := servicelocator.New("tracker", "yaml")
	result.AddConfigPath("$HOME/.tracker")
	result.AddConfigPath("/etc/tracker")

	result.RegisterConstructor("NewLogger", NewLogger)
	result.RegisterConstructor("NewUUID", NewUUID)
	result.RegisterConstructor("NewElasticClient", NewElasticClient)
	result.RegisterConstructor("NewDummyVisitRepository", NewDummyVisitRepository)
	result.RegisterConstructor("NewElasticVisitRepository", NewElasticVisitRepository)
	result.RegisterConstructor("NewMemoryVisitRepository", NewMemoryVisitRepository)
	result.RegisterConstructor("NewVisitManager", NewVisitManager)
	result.RegisterConstructor("NewDummyEventRepository", NewDummyEventRepository)
	result.RegisterConstructor("NewElasticEventRepository", NewElasticEventRepository)
	result.RegisterConstructor("NewMemoryEventRepository", NewMemoryEventRepository)
	result.RegisterConstructor("NewEventManager", NewEventManager)
	result.RegisterConstructor("NewDummyFlashRepository", NewDummyFlashRepository)
	result.RegisterConstructor("NewElasticFlashRepository", NewElasticFlashRepository)
	result.RegisterConstructor("NewFlashManager", NewFlashManager)
	result.RegisterConstructor("NewDummyMarkRepository", NewDummyMarkRepository)
	result.RegisterConstructor("NewMarkManager", NewMarkManager)
	result.RegisterConstructor("NewTrackManager", NewTrackManager)
	result.RegisterConstructor("NewVisitHandler", NewVisitHandler)
	result.RegisterConstructor("NewEventHandler", NewEventHandler)
	result.RegisterConstructor("NewFlashHandler", NewFlashHandler)
	result.RegisterConstructor("NewTrackHandler", NewTrackHandler)
	result.RegisterConstructor("NewMarkHandler", NewMarkHandler)
	result.RegisterConstructor("NewThriftService", NewThriftService)

	return result
}

func NewLogger(...interface{}) (interface{}, error) {
	return log.New(os.Stdout, "logger: ", log.Lshortfile), nil
}

func NewUUID(...interface{}) (interface{}, error) {
	return uuidDriver.New(), nil
}

func NewElasticClient(arguments ...interface{}) (interface{}, error) {
	options := []elasticClient.ClientOptionFunc{}

	if len(arguments) > 0 {
		for _, key := range arguments[0].([]interface{}) {
			options = append(options, elasticClient.SetURL(key.(string)))
		}
	}

	if len(arguments) == 2 {
		options = append(options, elasticClient.SetMaxRetries(arguments[1].(int)))
	}

	return elasticClient.NewClient(options...)
}

func NewDummyVisitRepository(arguments ...interface{}) (interface{}, error) {
	return dummy.NewVisitRepository(), nil
}

func NewElasticVisitRepository(arguments ...interface{}) (interface{}, error) {
	return elastic.NewVisitRepository(arguments[0].(*elasticClient.Client), arguments[1].(dao.UUIDProviderInterface))
}

func NewMemoryVisitRepository(arguments ...interface{}) (interface{}, error) {
	return memory.NewVisitRepository(arguments[0].(dao.VisitRepositoryInterface), arguments[1].(int))
}

func NewVisitManager(arguments ...interface{}) (interface{}, error) {
	return components.NewVisitManager(
		arguments[0].(dao.VisitRepositoryInterface),
		arguments[1].(dao.UUIDProviderInterface),
		arguments[2].(*log.Logger),
	), nil
}

func NewDummyEventRepository(arguments ...interface{}) (interface{}, error) {
	return dummy.NewEventRepository(), nil
}

func NewElasticEventRepository(arguments ...interface{}) (interface{}, error) {
	return elastic.NewEventRepository(arguments[0].(*elasticClient.Client), arguments[1].(dao.UUIDProviderInterface))
}

func NewMemoryEventRepository(arguments ...interface{}) (interface{}, error) {
	return memory.NewEventRepository(arguments[0].(dao.EventRepositoryInterface))
}

func NewEventManager(arguments ...interface{}) (interface{}, error) {
	return components.NewEventManager(
		arguments[0].(dao.EventRepositoryInterface),
		arguments[1].(dao.UUIDProviderInterface),
		arguments[2].(*log.Logger),
	), nil
}

func NewDummyFlashRepository(arguments ...interface{}) (interface{}, error) {
	return dummy.NewFlashRepository(), nil
}

func NewElasticFlashRepository(arguments ...interface{}) (interface{}, error) {
	return elastic.NewFlashRepository(arguments[0].(*elasticClient.Client), arguments[1].(dao.UUIDProviderInterface))
}

func NewFlashManager(arguments ...interface{}) (interface{}, error) {
	return components.NewFlashManager(
		arguments[0].(dao.FlashRepositoryInterface),
		arguments[1].(dao.UUIDProviderInterface),
		arguments[2].(*log.Logger),
	), nil
}

func NewTrackManager(arguments ...interface{}) (interface{}, error) {
	return components.NewTrackManager(
		arguments[0].(dao.VisitRepositoryInterface),
		arguments[1].(dao.EventRepositoryInterface),
		arguments[2].(dao.FlashRepositoryInterface),
		nil,
		arguments[4].(dao.UUIDProviderInterface),
		arguments[5].(*log.Logger),
	), nil
}

func NewDummyMarkRepository(arguments ...interface{}) (interface{}, error) {
	return dummy.NewMarkRepository(), nil
}

func NewMarkManager(arguments ...interface{}) (interface{}, error) {
	return components.NewMarkManager(
		arguments[0].(dao.MarkRepositoryInterface),
		arguments[1].(dao.UUIDProviderInterface),
		arguments[2].(*log.Logger),
	), nil
}

func NewVisitHandler(arguments ...interface{}) (interface{}, error) {
	return handlers.NewVisitHandler(
		arguments[0].(*components.VisitManager),
		arguments[1].(dao.UUIDProviderInterface),
	), nil
}

func NewEventHandler(arguments ...interface{}) (interface{}, error) {
	return handlers.NewEventHandler(
		arguments[0].(*components.EventManager),
		arguments[1].(dao.UUIDProviderInterface),
	), nil
}

func NewFlashHandler(arguments ...interface{}) (interface{}, error) {
	return handlers.NewFlashHandler(
		arguments[0].(*components.FlashManager),
		arguments[1].(dao.UUIDProviderInterface),
	), nil
}

func NewTrackHandler(arguments ...interface{}) (interface{}, error) {
	return handlers.NewTrackHandler(
		arguments[0].(*components.TrackManager),
		arguments[1].(dao.UUIDProviderInterface),
	), nil
}

func NewMarkHandler(arguments ...interface{}) (interface{}, error) {
	return handlers.NewMarkHandler(
		arguments[0].(*components.MarkManager),
		arguments[1].(dao.UUIDProviderInterface),
	), nil
}

func NewThriftService(arguments ...interface{}) (interface{}, error) {
	processor := thrift.NewTMultiplexedProcessor()
	processor.RegisterProcessor("visit", generated.NewVisitServiceProcessor(arguments[0].(*handlers.VisitHandler)))
	processor.RegisterProcessor("event", generated.NewEventServiceProcessor(arguments[1].(*handlers.EventHandler)))
	processor.RegisterProcessor("flash", generated.NewFlashServiceProcessor(arguments[2].(*handlers.FlashHandler)))
	processor.RegisterProcessor("track", generated.NewTrackServiceProcessor(arguments[3].(*handlers.TrackHandler)))
	processor.RegisterProcessor("mark", generated.NewMarkServiceProcessor(arguments[4].(*handlers.MarkHandler)))

	transport, err := thrift.NewTServerSocket(arguments[5].(string))
	if err != nil {
		panic(err.Error())
	}

	size, _ := strconv.Atoi(arguments[6].(string))

	transportFactory := thrift.NewTBufferedTransportFactory(size)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	return thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory), nil
}
