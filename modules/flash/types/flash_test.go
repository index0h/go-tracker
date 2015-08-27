package entities

import (
	"testing"
	"time"

	"github.com/index0h/go-tracker/dao/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFlash_NewFlash_EmptyFlashID(t *testing.T) {
	flash, err := NewFlash([16]byte{}, time.Now().Unix(), &Visit{}, &Event{})

	assert.Nil(t, flash)
	assert.NotNil(t, err)
}

func TestFlash_NewFlash_EmptyEvent(t *testing.T) {
	flash, err := NewFlash(uuid.New().Generate(), time.Now().Unix(), &Visit{}, nil)

	assert.Nil(t, flash)
	assert.NotNil(t, err)
}

func TestFlash_NewFlash_EmptyVisit(t *testing.T) {
	flash, err := NewFlash(uuid.New().Generate(), time.Now().Unix(), nil, &Event{})

	assert.Nil(t, flash)
	assert.NotNil(t, err)
}
