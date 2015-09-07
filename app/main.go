package main

import (
	"github.com/index0h/go-tracker/modules/event"
	"github.com/index0h/go-tracker/modules/flash"
	"github.com/index0h/go-tracker/modules/track"
	"github.com/index0h/go-tracker/modules/visit"

	eventDummy "github.com/index0h/go-tracker/modules/event/dao/dummy"
	flashDummy "github.com/index0h/go-tracker/modules/flash/dao/dummy"
	visitDummy "github.com/index0h/go-tracker/modules/visit/dao/dummy"

	eventElastic "github.com/index0h/go-tracker/modules/event/dao/elastic"
	flashElastic "github.com/index0h/go-tracker/modules/flash/dao/elastic"
	visitElastic "github.com/index0h/go-tracker/modules/visit/dao/elastic"

	eventMemory "github.com/index0h/go-tracker/modules/event/dao/memory"
	visitMemory "github.com/index0h/go-tracker/modules/visit/dao/memory"

	"github.com/index0h/go-tracker/app/generated"
	"github.com/index0h/go-tracker/app/handlers"
	"github.com/index0h/go-tracker/share/uuid"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Sirupsen/logrus"
	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"
	"os"
	"strconv"
)

const Version = "0.0.1"

type Config struct {
	Port     uint
	Host     string
	LogLevel string

	ElasticConnections struct {
		Host       string
		MaxRetries int
	}

	Visit struct {
		Elastic bool
		Memory  struct {
			Use       bool
			CacheSize int
		}
	}

	Event struct {
		Elastic bool
		Memory  struct {
			Use bool
		}
	}

	Flash struct {
		Elastic bool
	}
}

func main() {
	var (
		config     Config
		configPath string
		err        error
	)

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	configLoader := viper.New()
	configLoader.SetConfigType("yml")
	configLoader.SetConfigName("tracker")

	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run tracker service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := configLoader.ReadInConfig(); err != nil {
				logger.Fatal(err)
			}

			if err := configLoader.Marshal(&config); err != nil {
				logger.Fatal(err)
			}

			if logger.Level, err = logrus.ParseLevel(config.LogLevel); err != nil {
				logger.Panic(err)
			}

			Run(&config, logger)
		},
	}

	runCommand.Flags().StringVarP(&configPath, "config", "c", "", "alternative config path")
	if configPath != "" {
		configLoader.AddConfigPath(configPath)
	} else {
		currentPath, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
		}

		configLoader.AddConfigPath(currentPath)
		configLoader.AddConfigPath("/etc/tracker/")
		configLoader.AddConfigPath("$HOME/.tracker")
	}

	runCommand.Flags().StringVarP(&config.LogLevel, "log", "l", "warning", "log level")
	viper.BindPFlag("Log", runCommand.Flags().Lookup("log"))

	runCommand.Flags().StringVar(&config.Host, "host", "localhost", "tracker service host")
	viper.BindPFlag("Host", runCommand.Flags().Lookup("host"))

	runCommand.Flags().UintVarP(&config.Port, "port", "p", 9898, "tracker service port")
	viper.BindPFlag("Port", runCommand.Flags().Lookup("port"))

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Print tracker version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}

	rootCommand := &cobra.Command{Use: "app"}
	rootCommand.AddCommand(runCommand, versionCommand)
	rootCommand.Execute()
}

func Run(config *Config, logger *logrus.Logger) {
	var (
		visitRepository visit.RepositoryInterface
		eventRepository event.RepositoryInterface
		flashRepository flash.RepositoryInterface
		elasticClient   *elastic.Client
		err             error
	)

	uuid := uuid.New()

	if config.ElasticConnections.Host != "" {
		logger.Info("create elastic client")

		elasticClient, err = elastic.NewClient(
			elastic.SetURL(config.ElasticConnections.Host),
			elastic.SetMaxRetries(config.ElasticConnections.MaxRetries),
		)
		if err != nil {
			logger.Fatal(err)
		}
	}

	visitRepository = visitDummy.NewRepository()

	if config.Visit.Elastic {
		logger.Info("create elastic visit repository")

		visitRepository, err = visitElastic.NewRepository(elasticClient, uuid)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if config.Visit.Memory.Use {
		logger.Info("create memory visit repository")

		visitRepository, err = visitMemory.NewRepository(visitRepository, config.Visit.Memory.CacheSize)
		if err != nil {
			logger.Fatal(err)
		}
	}

	visitManager := visit.NewManager(visitRepository, uuid, logger)

	eventRepository = eventDummy.NewRepository()

	if config.Event.Elastic {
		logger.Info("create elastic event repository")

		eventRepository, err = eventElastic.NewRepository(elasticClient, uuid)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if config.Event.Memory.Use {
		logger.Info("create memory event repository")

		eventRepository, err = eventMemory.NewRepository(eventRepository)
		if err != nil {
			logger.Fatal(err)
		}
	}

	eventManager := event.NewManager(eventRepository, uuid, logger)

	flashRepository = flashDummy.NewRepository()

	if config.Flash.Elastic {
		logger.Info("create elastic flash repository")

		flashRepository, err = flashElastic.NewRepository(elasticClient, uuid)
		if err != nil {
			logger.Fatal(err)
		}
	}

	flashManager := flash.NewManager(flashRepository, uuid, logger)

	trackManager := track.NewManager(visitManager, eventManager, flashManager, nil, uuid, logger)

	logger.Info("init thrift")

	processor := thrift.NewTMultiplexedProcessor()

	visitHandler := handlers.NewVisitHandler(visitManager, uuid)
	eventHandler := handlers.NewEventHandler(eventManager, uuid)
	flashHandler := handlers.NewFlashHandler(flashManager, uuid)
	trackHandler := handlers.NewTrackHandler(trackManager, uuid)

	processor.RegisterProcessor("visit", generated.NewVisitServiceProcessor(visitHandler))
	processor.RegisterProcessor("event", generated.NewEventServiceProcessor(eventHandler))
	processor.RegisterProcessor("flash", generated.NewFlashServiceProcessor(flashHandler))
	processor.RegisterProcessor("track", generated.NewTrackServiceProcessor(trackHandler))

	transport, err := thrift.NewTServerSocket(config.Host + ":" + strconv.Itoa(int(config.Port)))
	if err != nil {
		logger.Fatal(err)
	}

	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	logger.Info("run server on: tcp://" + config.Host + ":" + strconv.Itoa(int(config.Port)))

	server.Serve()
}
