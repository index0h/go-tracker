package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao/elastic"
	"github.com/index0h/go-tracker/dao/memory"
	uuidDriver "github.com/index0h/go-tracker/dao/uuid"
	elasticClient "github.com/olivere/elastic"
	"log"
	"os"
	"github.com/index0h/go-tracker/app/handlers"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	uuid := uuidDriver.New()

	elasticClient, err := elasticClient.NewClient(elasticClient.SetURL("127.0.0.1"), elasticClient.SetMaxRetries(10))
	if err != nil {
		panic(err.Error())
	}

	visitElasticRepository, err := elastic.NewVisitRepository(elasticClient, uuid)
	if err != nil {
		panic(err.Error())
	}

	visitRepository, err := memory.NewVisitRepository(visitElasticRepository, 2000)
	if err != nil {
		panic(err.Error())
	}

	visitManager, err := components.NewVisitManager(visitRepository, uuid, logger)
	if err != nil {
		panic(err.Error())
	}

	eventElasticRepository, err := elastic.NewEventRepository(elasticClient, uuid)
	if err != nil {
		panic(err.Error())
	}

	eventRepository, err := memory.NewEventRepository(eventElasticRepository)
	if err != nil {
		panic(err.Error())
	}

	eventManager, err := components.NewEventManager(eventRepository, uuid, logger)
	if err != nil {
		panic(err.Error())
	}

	flashRepository, err := elastic.NewFlashRepository(elasticClient, uuid)
	if err != nil {
		panic(err.Error())
	}

	flashManager, err := components.NewVisitManager(flashRepository, uuid, logger)
	if err != nil {
		panic(err.Error())
	}

	trackManager, err := components.NewTrackManager(visitRepository, eventRepository, flashRepository, nil, uuid, logger)
	if err != nil {
		panic(err.Error())
	}

	markManager, err := components.MarkManager(nil, uuid, logger)
	if err != nil {
		panic(err.Error())
	}

	visitHandler := handlers.NewVisitHandler(visitManager, uuid)
	eventHandler := handlers.NewEventHandler(eventManager, uuid)
	flashHandler := handlers.NewFlashHandler(flashManager, uuid)
	trackHandler := handlers.NewTrackHandler(trackManager, uuid)
	markHandler := handlers.NewMarkHandler(markManager, uuid)

	processor := thrift.NewTMultiplexedProcessor()
	processor.RegisterProcessor("visit", generated.NewVisitServiceProcessor(visitHandler))
	processor.RegisterProcessor("event", generated.NewVisitServiceProcessor(eventHandler))
	processor.RegisterProcessor("flash", generated.NewVisitServiceProcessor(flashHandler))
	processor.RegisterProcessor("track", generated.NewVisitServiceProcessor(trackHandler))
	processor.RegisterProcessor("mark", generated.NewVisitServiceProcessor(markHandler))

	transport, err := thrift.NewTServerSocket("127.0.0.1:9090")
	if err != nil {
		panic(err.Error())
	}

	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	server.Serve()
}
