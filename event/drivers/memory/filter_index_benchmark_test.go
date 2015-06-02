package memory

import (
	"math/rand"
	"testing"
	"time"

	eventEntities "github.com/index0h/go-tracker/event/entities"
	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	visitEntities "github.com/index0h/go-tracker/visit/entities"
)

func BenchmarkFilterIndexFindAllByVisit3(b *testing.B) {
	filterIndexFindAllByVisit(3, b)
}

func BenchmarkFilterIndexFindAllByVisit5(b *testing.B) {
	filterIndexFindAllByVisit(5, b)
}

func BenchmarkFilterIndexFindAllByVisit10(b *testing.B) {
	filterIndexFindAllByVisit(10, b)
}

func BenchmarkFilterIndexFindAllByVisit15(b *testing.B) {
	filterIndexFindAllByVisit(15, b)
}

var filterIndexSymbols = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func filterIndexFindAllByVisit(countKeys uint, b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().UTC().UnixNano())

	events := filterIndexGenerateEvents(uint(b.N), countKeys)
	visits := filterIndexGenerateVisits(uint(b.N), countKeys)

	index := NewFilterIndex()
	index.Refresh(events)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = index.FindAllByVisit(visits[i])
	}
}

func filterIndexGenerateVisits(countVisits uint, countData uint) []*visitEntities.Visit {
	result := make([]*visitEntities.Visit, countVisits)

	uuid := uuidDriver.UUID{}

	for i := uint(0); i < countVisits; i++ {
		data := filterIndexGenerateKeyValue(countData)

		result[i], _ = visitEntities.NewVisit(uuid.Generate(), int64(0), uuid.Generate(), "", data, []string{})
	}

	return result
}

func filterIndexGenerateEvents(count uint, countData uint) []*eventEntities.Event {
	result := make([]*eventEntities.Event, count)

	uuid := uuidDriver.UUID{}

	for i := uint(0); i < count; i++ {
		filters := filterIndexGenerateKeyValue(countData)

		result[i], _ = eventEntities.NewEvent(uuid.Generate(), true, map[string]string{}, filters)
	}

	return result
}

func filterIndexGenerateKeyValue(count uint) (result map[string]string) {
	result = make(map[string]string, count)

	for i := uint(0); i < count; i++ {
		result[filterIndexGenerateString()] = filterIndexGenerateString()
	}

	return result
}

func filterIndexGenerateString() string {
	count := 5
	result := make([]byte, count)

	var number int

	for i := 0; i < count; i++ {
		number = rand.Intn(len(filterIndexSymbols) - 1)

		result[i] = filterIndexSymbols[number]
	}

	return string(result)
}
