package memory

import (
	"math/rand"
	"testing"
	"time"

	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	eventEntities "github.com/index0h/go-tracker/event/entities"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

func BenchmarkFindAllByVisit3(b *testing.B) {
	indexFindAllByVisit(3, b)
}

func BenchmarkFindAllByVisit5(b *testing.B) {
	indexFindAllByVisit(5, b)
}

func BenchmarkFindAllByVisit10(b *testing.B) {
	indexFindAllByVisit(10, b)
}

func BenchmarkFindAllByVisit15(b *testing.B) {
	indexFindAllByVisit(15, b)
}

var symbols = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func indexFindAllByVisit(countKeys uint, b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().UTC().UnixNano())

	events := generateEvents(uint(b.N), countKeys)
	visits := generateVisits(uint(b.N), countKeys)

	index := NewFilterIndex()
	index.InsertAll(events)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = index.FindAllByVisit(visits[i])
	}
}

func generateVisits(countVisits uint, countData uint) ([]*visitEntities.Visit) {
	result := make([]*visitEntities.Visit, countVisits)

	uuid := uuidDriver.UUID{}

	for i := uint(0); i < countVisits; i++ {
		data := generateKeyValue(countData)

		result[i], _ = visitEntities.NewVisit(uuid.Generate(), int64(0), uuid.Generate(), "", data, []string{})
	}

	return result
}

func generateEvents(count uint, countData uint) ([]*eventEntities.Event) {
	result := make([]*eventEntities.Event, count)

	uuid := uuidDriver.UUID{}

	for i := uint(0); i < count; i++ {
		filters := generateKeyValue(countData)

		result[i], _ = eventEntities.NewEvent(uuid.Generate(), true, map[string]string{}, filters)
	}

	return result
}

func generateKeyValue(count uint) (result map[string]string) {
	result = make(map[string]string, count)

	for i := uint(0); i < count; i++ {
		result[generateString()] = generateString()
	}

	return result
}

func generateString() string {
	count := 5
	result := make([]byte, count)

	var number int

	for i := 0; i < count; i++ {
		number = rand.Intn(len(symbols) - 1)

		result[i] = symbols[number]
	}

	return string(result)
}

