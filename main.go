package main

import (
	"github.com/index0h/go-tracker/uuid/drivers/uuidDriver"
	"github.com/index0h/go-tracker/visit"
	"github.com/index0h/go-tracker/visit/drivers/elasticDriver"
	"github.com/index0h/go-tracker/visit/drivers/memoryDriver"
	"github.com/olivere/elastic"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	uuid := new(uuidDriver.UUID)
	client, _ := elastic.NewClient()
	repository := elasticDriver.NewRepository(client, uuid)
	memoryRepository := memoryDriver.NewRepository(repository, 100)
	manager := visit.NewManager(memoryRepository, uuid, logger)

	visit, err := manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})
	_, _ = manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})

	logger.Println(visit)
	logger.Println(err)
}