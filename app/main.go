package main

import (
	thriftGenerated "./tracker"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/index0h/go-tracker/components"
	"github.com/index0h/go-tracker/dao"
	"github.com/index0h/go-tracker/dao/dummy"
	"github.com/index0h/go-tracker/dao/elastic"
	"github.com/index0h/go-tracker/dao/memory"
	"github.com/index0h/go-tracker/dao/processors"
	uuidDriver "github.com/index0h/go-tracker/dao/uuid"
	elasticClient "github.com/olivere/elastic"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	configStorage  *viper.Viper
	loggerStorage  *log.Logger
	uuidStorage    dao.UUIDProviderInterface
	elasticStorage map[string]*elasticClient.Client
)

func main() {
	server := getThriftServer()
	server.Serve()
}

func getVisitManager() *components.VisitManager {
	var repository dao.VisitRepositoryInterface

	switch getConfig().GetString("visit_manager.repository.type") {
	case "dummy":
		repository = getVisitRepositoryDummy()
	case "elasticsearch":
		repository = getVisitRepositoryElastic("visit_manager.repository")
	case "memory":
		repository = getVisitRepositoryMemory("visit_manager.repository")
	default:
		panic("Invalid visit repository type")
	}

	return components.NewVisitManager(repository, getUUID(), getLogger())
}

func getEventManager() *components.EventManager {
	var repository dao.EventRepositoryInterface

	switch getConfig().GetString("event_manager.repository.type") {
	case "dummy":
		repository = getEventRepositoryDummy()
	case "elasticsearch":
		repository = getEventRepositoryElastic("event_manager.repository")
	case "memory":
		repository = getEventRepositoryMemory("event_manager.repository")
	default:
		panic("Invalid event repository type")
	}

	return components.NewEventManager(repository, getUUID(), getLogger())
}

func getFlashManager() *components.FlashManager {
	var repository dao.FlashRepositoryInterface

	switch getConfig().GetString("flash_manager.repository.type") {
	case "dummy":
		repository = getFlashRepositoryDummy()
	case "elasticsearch":
		repository = getFlashRepositoryElastic("flash_manager.repository")
	default:
		panic("Invalid flash repository type")
	}

	return components.NewFlashManager(repository, getUUID(), getLogger())
}

func getTrackerManager() *components.TrackerManager {
	return components.NewTrackerManager(
		getVisitManager(),
		getEventManager(),
		getFlashManager(),
		getProcessors(),
		getUUID(),
		getLogger(),
	)
}

func getThriftServer() *thrift.TSimpleServer {
	return thrift.NewTSimpleServer4(
		getThriftProcessor(),
		getThriftTransport(),
		getThriftTransportFactory(),
		getThriftProtocolFactory(),
	)
}

func getConfig() *viper.Viper {
	if configStorage != nil {
		return configStorage
	}

	configStorage = viper.New()
	configStorage.AddConfigPath("config")
	configStorage.AddConfigPath("gotracker")
	configStorage.AddConfigPath("/etc/gotracker")
	configStorage.AddConfigPath("$HOME/.gotracker")

	if err := configStorage.ReadInConfig(); err != nil {
		panic("Fatal error config file")
	}

	return configStorage
}

func getLogger() *log.Logger {
	if loggerStorage == nil {
		loggerStorage = log.New(os.Stdout, "logger: ", log.Lshortfile)
	}

	return loggerStorage
}

func getUUID() dao.UUIDProviderInterface {
	if uuidStorage == nil {
		uuidStorage = uuidDriver.New()
	}

	return uuidStorage
}

func getElasticClient(name string) *elasticClient.Client {
	if len(elasticStorage) == 0 {
		config := getConfig().GetStringMapString("elasticsearch")

		elasticStorage = make(map[string]*elasticClient.Client, len(config))

		for elasticName := range config {
			configPrefix := "elasticsearch." + elasticName

			hosts := elasticClient.SetURL(getConfig().GetStringSlice(configPrefix + ".hosts")...)
			maxRetries := elasticClient.SetMaxRetries(getConfig().GetInt(configPrefix + ".max_retries"))

			client, err := elasticClient.NewClient(hosts, maxRetries)
			if err != nil {
				panic("Invalid elastic connection")
			}

			elasticStorage[elasticName] = client
		}
	}

	result, ok := elasticStorage[name]
	if !ok {
		panic("ElasticSearch client not found")
	}

	return result
}

func getVisitRepositoryDummy() *dummy.VisitRepository {
	return new(dummy.VisitRepository)
}

func getVisitRepositoryElastic(configPath string) *elastic.VisitRepository {
	client := getElasticClient(getConfig().GetString(configPath + ".client"))

	repository, err := elastic.NewVisitRepository(client, getUUID())
	if err != nil {
		panic(err.Error())
	}

	return repository
}

func getVisitRepositoryMemory(configPath string) *memory.VisitRepository {
	var nested dao.VisitRepositoryInterface

	switch getConfig().GetString(configPath + ".nested.type") {
	case "dummy":
		nested = getVisitRepositoryDummy()
	case "elasticsearch":
		nested = getVisitRepositoryElastic(configPath + ".nested")
	default:
		panic("Invalid nested visit repository type")
	}

	maxEntries := getConfig().GetInt(configPath + ".max_entities")

	repository, err := memory.NewVisitRepository(nested, maxEntries)
	if err != nil {
		panic(err.Error())
	}

	return repository
}

func getEventRepositoryDummy() *dummy.EventRepository {
	return new(dummy.EventRepository)
}

func getEventRepositoryElastic(configPath string) *elastic.EventRepository {
	client := getElasticClient(getConfig().GetString(configPath + ".client"))

	repository, err := elastic.NewEventRepository(client, getUUID())
	if err != nil {
		panic(err.Error())
	}

	return repository
}

func getEventRepositoryMemory(configPath string) *memory.EventRepository {
	var nested dao.EventRepositoryInterface

	switch getConfig().GetString(configPath + ".nested.type") {
	case "dummy":
		nested = getEventRepositoryDummy()
	case "elasticsearch":
		nested = getEventRepositoryElastic(configPath + ".nested")
	default:
		panic("Invalid nested event repository type")
	}

	repository, err := memory.NewEventRepository(nested)
	if err != nil {
		panic(err.Error())
	}

	return repository
}

func getFlashRepositoryDummy() *dummy.FlashRepository {
	return new(dummy.FlashRepository)
}

func getFlashRepositoryElastic(configPath string) *elastic.FlashRepository {
	client := getElasticClient(getConfig().GetString(configPath + ".client"))

	repository, err := elastic.NewFlashRepository(client, getUUID())
	if err != nil {
		panic(err.Error())
	}

	return repository
}

func getProcessors() []dao.ProcessorInterface {
	config := getConfig().GetStringMapString("tracker_manager.processors")

	result := make([]dao.ProcessorInterface, len(config))

	var i = 0
	for processorType := range config {
		priority := getConfig().GetInt("tracker_manager.processors." + processorType + ".priority")

		switch processorType {
		case "dummy":
			result[i] = processors.NewDummy(priority)
		default:
			panic("Invalid processor type")
		}
	}

	return result
}

func getThriftHandler() *ThriftHandler {
	return &ThriftHandler{
		visitManager:   getVisitManager(),
		eventManager:   getEventManager(),
		flashManager:   getFlashManager(),
		trackerManager: getTrackerManager(),
		uuid:           getUUID(),
		logger:         getLogger(),
	}
}

func getThriftProcessor() thrift.TProcessor {
	return thriftGenerated.NewTrackerServiceProcessor(getThriftHandler())
}

func getThriftTransport() *thrift.TServerSocket {
	result, err := thrift.NewTServerSocket(getConfig().GetString("thrift.addr"))
	if err != nil {
		panic(err.Error())
	}

	return result
}

func getThriftTransportFactory() thrift.TTransportFactory {
	transportFactory := thrift.NewTBufferedTransportFactory(getConfig().GetInt("thrift.buffer_size"))

	return thrift.NewTFramedTransportFactory(transportFactory)
}

func getThriftProtocolFactory() thrift.TProtocolFactory {
	switch getConfig().GetString("thrift.protocol") {
	case "compact":
		return thrift.NewTCompactProtocolFactory()
	case "simplejson":
		return thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		return thrift.NewTJSONProtocolFactory()
	case "binary":
		return thrift.NewTBinaryProtocolFactoryDefault()
	default:
		panic("Invalid protocol specified")
	}
}
