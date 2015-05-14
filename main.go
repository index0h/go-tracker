package main

import (
	"github.com/index0h/go-tracker/uuid/drivers/satori/uuidDriver"
	"github.com/index0h/go-tracker/drivers/olivere/elastic/visitDriver"
	"github.com/index0h/go-tracker/visit"
	"github.com/olivere/elastic"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	uuid := uuidDriver.Uuid{}
	client, _ := elastic.NewClient()
	repository := visitDriver.NewRepository(client, uuid)
	manager := visit.NewManager(repository, uuid, logger)

	visit, err := manager.Track(uuid.Generate(), "", map[string]string{"A": "B"})

	logger.Println(visit)
	logger.Println(err)
}