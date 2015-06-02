package memory

import (
	"testing"

	eventPackage "github.com/index0h/go-tracker/event"
	//	eventEntities "github.com/index0h/go-tracker/event/entities"
	//	interfaceUUID "github.com/index0h/go-tracker/uuid"
	//	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	//	visitEntities "github.com/index0h/go-tracker/visit/entities"
	//	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	func(event eventPackage.Repository) {}(&Repository{})
}
